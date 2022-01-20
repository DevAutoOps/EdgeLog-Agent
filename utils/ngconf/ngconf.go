package ngconf

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// Pre-defines some errors.
var (
	ErrRootDirective      = fmt.Errorf("root node must have the directive")
	ErrNonRootNoDirective = fmt.Errorf("non-root node have no directive")
)

type nodes []*Node

func (ns nodes) Len() int      { return len(ns) }
func (ns nodes) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes) Less(i, j int) bool {
	if ns[j] == nil {
		return true
	}
	return false
}

type nodeStack struct {
	nodes []*Node
}

func (ns *nodeStack) Push(node *Node) { ns.nodes = append(ns.nodes, node) }
func (ns *nodeStack) Pop() *Node {
	_len := len(ns.nodes) - 1
	node := ns.nodes[_len]
	ns.nodes = ns.nodes[:_len]
	return node
}

type MapNodes map[string]*Node

// Node represents a set of the key-value configurations of Nginx.
type Node struct {
	Root      bool
	Directive string
	Args      []string
	Children  []*Node
}

func newNode(stmts []string, root ...bool) (*Node, error) {
	var isRoot bool
	if len(root) > 0 && root[0] {
		isRoot = true
	}

	var directive string
	var args []string
	switch len(stmts) {
	case 0:

	case 1:
		directive = stmts[0]
	default:
		directive = stmts[0]
		args = stmts[1:]
	}

	if directive == "" && !isRoot {
		return nil, ErrNonRootNoDirective
	}
	if directive != "" && isRoot {
		return nil, ErrRootDirective
	}

	return &Node{Directive: directive, Args: args, Root: isRoot}, nil
}

func (n *Node) appendChild(node *Node) { n.Children = append(n.Children, node) }

// String is equal to n.Dump(0).
func (n *Node) String() string {
	return n.Dump(0)
}

// WriteTo implements the interface WriterTo to write the configuration to file.
func (n *Node) WriteTo(w io.Writer) (int64, error) {
	m, err := io.WriteString(w, n.String())
	return int64(m), err
}

// WriteToFile Write to file
func (n *Node) WriteToFile(filePath string) (int64, error) {
	var f *os.File
	var err error
	if checkFileIsExist(filePath) { //If the file exists
		f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666) //open file
	} else {
		f, err = os.Create(filePath) //create file
	}
	if err != nil {
		return 0, err
	}
	return n.WriteTo(f)
}

func (n *Node) getChildren(indent int, ctx *nodeDumpCtx) string {
	ss := make([]string, len(n.Children))
	for i, node := range n.Children {
		if i == 0 {
			ctx.LastBlockEnd = false
		} else {
			ctx.BlockStart = false
		}
		ss[i] = node.dump(indent, ctx)
	}
	return strings.Join(ss, "\n")
}

// Dump converts the Node to string.
func (n *Node) Dump(indent int) string {
	return n.dump(indent, &nodeDumpCtx{FirstBlock: true})
}

type nodeDumpCtx struct {
	HasComment   bool
	FirstBlock   bool
	BlockStart   bool
	LastBlockEnd bool
}

func (n *Node) dump(indent int, ctx *nodeDumpCtx) string {
	var prefix, spaces string
	for i := indent; i > 0; i-- {
		spaces += "    "
	}

	lastComment := ctx.HasComment
	ctx.HasComment = strings.HasPrefix(n.Directive, "#")
	if ctx.HasComment && !lastComment {
		if !ctx.BlockStart {
			prefix = "\n"
		}
	} else if ctx.LastBlockEnd {
		prefix = "\n"
	}

	if len(n.Children) == 0 && n.Directive != "location" {
		return fmt.Sprintf("%s%s%s %s;", prefix, spaces, n.Directive, strings.Join(n.Args, " "))
	} else if n.Root {
		return n.getChildren(indent, ctx)
	} else if args := strings.Join(n.Args, " "); args != "" {
		if ctx.FirstBlock {
			ctx.FirstBlock = false
			if prefix == "" {
				prefix = "\n"
			}
		}
		ctx.BlockStart = true
		s := fmt.Sprintf("%s%s%s %s {\n%s\n%s}", prefix, spaces, n.Directive, args,
			n.getChildren(indent+1, ctx), spaces)
		ctx.LastBlockEnd = true
		return s
	} else {
		if ctx.FirstBlock {
			ctx.FirstBlock = false
			if prefix == "" {
				prefix = "\n"
			}
		}
		ctx.BlockStart = true
		s := fmt.Sprintf("%s%s%s {\n%s\n%s}", prefix, spaces, n.Directive,
			n.getChildren(indent+1, ctx), spaces)
		ctx.LastBlockEnd = true
		return s
	}
}

// Get returns the child node by the given directive with the args.
func (n *Node) Get(directive string, args ...string) []*Node {
	results := make([]*Node, 0, 1)
	for _, child := range n.Children {
		if strings.ToLower(child.Directive) == strings.ToLower(directive) {
			results = append(results, child)
		}
	}

	_len := len(results)
	if _len == 0 {
		return nil
	}

	if argslen := len(args); argslen > 0 {
		var count int
		for i, node := range results {
			if len(node.Args) < argslen {
				results[i] = nil
				count++
			}
		}
		if count == _len {
			return nil
		} else if count > 0 {
			sort.Sort(nodes(results))
			results = results[:_len-count]
		}

		oldResults := results
		results = make([]*Node, 0, len(oldResults))
		for _, node := range oldResults {
			ok := true
			for i, arg := range args {
				if arg != node.Args[i] {
					ok = false
					break
				}
			}
			if ok {
				results = append(results, node)
			}
		}
	}

	return results
}

func (n *Node) GetDeepByDirective(directive string) []*Node {
	results := make([]*Node, 0, 1)
	for _, child := range n.Children {
		if strings.ToLower(child.Directive) == strings.ToLower(directive) {
			results = append(results, child)
			continue
		}
		if child.Children != nil && len(child.Children) > 0 {
			childNodes := child.GetDeepByDirective(directive)
			if len(childNodes) > 0 {
				results = append(results, childNodes...)
			}
		}
	}

	return results
}

func (n *Node) GetDeepByDirectiveAndArgs(directive string, args string) []*Node {
	results := make([]*Node, 0, 1)
	for _, child := range n.Children {
		argStr := ""
		if child.Args != nil && len(child.Args) >= 0 {
			argStr = strings.Join(child.Args, " ")
		}
		if strings.ToLower(child.Directive) == strings.ToLower(directive) && argStr == args {
			results = append(results, child)
			continue
		}
		if child.Children != nil && len(child.Children) > 0 {
			childNodes := child.GetDeepByDirectiveAndArgs(directive, args)
			if len(childNodes) > 0 {
				results = append(results, childNodes...)
			}
		}
	}

	return results
}

func (n *Node) GetDeep(directive string, exclude string) []*Node {
	results := make([]*Node, 0, 1)
	for _, child := range n.Children {
		if len(exclude) > 0 && strings.ToLower(child.Directive) == strings.ToLower(exclude) {
			continue
		}
		if strings.ToLower(child.Directive) == strings.ToLower(directive) {
			results = append(results, child)
			continue
		}
		if child.Children != nil && len(child.Children) > 0 {
			childNodes := child.GetDeep(directive, exclude)
			if len(childNodes) > 0 {
				results = append(results, childNodes...)
			}
		}
	}

	return results
}

// GetParent Get parent node
func (n *Node) GetParent(childNode *Node, rootNode *Node) (parentNode *Node) {
	if len(rootNode.Children) == 0 {
		return nil
	}
	for _, child := range rootNode.Children {
		if child == childNode {
			return rootNode
		}
		if child.Children != nil && len(child.Children) > 0 {
			parentNode = child.GetParent(childNode, child)
		}
	}
	return parentNode
}

// Add adds and returns the child node with the directive and the args.
//
// If the child node has existed, it will be ignored and return the first old.
func (n *Node) Add(directive string, args ...string) *Node {
	if nodes := n.Get(directive, args...); len(nodes) > 0 {
		return nodes[0]
	}

	node := &Node{Directive: directive, Args: args}
	n.appendChild(node)
	return node
}

func (n *Node) AddWithChildren(directive string, children []*Node, args ...string) *Node {
	if nodes := n.Get(directive, args...); len(nodes) > 0 {
		return nodes[0]
	}

	node := &Node{Directive: directive, Args: args, Children: children}
	n.appendChild(node)
	return node
}

// Del deletes the child node by the directive with the args.
//
// If args is nil, it will delete all the child nodes.
func (n *Node) Del(directive string, args ...string) {
	_len := len(n.Children)
	if _len == 0 {
		return
	}

	var count int
	for i, child := range n.Children {
		if strings.ToLower(child.Directive) == strings.ToLower(directive) {
			if _len := len(args); _len == 0 {
				n.Children[i] = nil
				count++
			} else if _len <= len(child.Args) {
				ok := true
				for j, arg := range args {
					if arg != child.Args[j] {
						ok = false
						break
					}
				}
				if ok {
					n.Children[i] = nil
					count++
				}
			}
		}
	}

	if count > 0 {
		sort.Sort(nodes(n.Children))
		n.Children = n.Children[:_len-count]
	}
}

// DelNode Deletes child nodes from the specified node
func (n *Node) DelNode(delNode *Node, rootNode *Node) {
	if len(rootNode.Children) == 0 {
		return
	}
	for i, child := range rootNode.Children {
		if child == delNode {
			if i == 0 {
				rootNode.Children = rootNode.Children[1:]
			} else if i == len(rootNode.Children)-1 {
				rootNode.Children = rootNode.Children[:i]
			} else {
				rootNode.Children = append(rootNode.Children[:i], rootNode.Children[i+1:]...)
			}
			break
		}
		if child.Children != nil && len(child.Children) > 0 {
			child.DelNode(delNode, child)
		}
	}
}

// Decode decodes the string s to Node.
func Decode(s string) (*Node, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ngconf decode runtime panic caught: %v\n", err)
			return
		}
	}()

	var err error
	var node *Node
	stack := &nodeStack{}
	currentWord := []rune{}
	currentStmt := []string{}
	currentBlock, _ := newNode(nil, true)

	isSingleQuotation := false
	isComment := false
	for _, char := range s {
		if char == '\n' {
			isComment = false
		}
		if isComment {
			continue
		}
		switch char {
		case '#':
			isComment = true
		case '\'':
			isSingleQuotation = !isSingleQuotation
			currentWord = append(currentWord, char)
		case '{':
			// Put the current block on the stack, start a new block.
			// Also, if we are in a word, "finish" that off, and end
			// the current statement.
			if isSingleQuotation {
				currentWord = append(currentWord, char)
				continue
			}
			stack.Push(currentBlock)
			if len(currentWord) > 0 {
				currentStmt = append(currentStmt, string(currentWord))
				currentWord = nil
			}
			if currentBlock, err = newNode(currentStmt); err != nil {
				return nil, err
			}
			currentStmt = nil
		case '}':
			// Finalize the current block, pull the previous (outer) block off
			// of the stack, and add the inner block to the previous block's
			// map of blocks.
			if isSingleQuotation {
				currentWord = append(currentWord, char)
				continue
			}
			innerBlock := currentBlock
			currentBlock = stack.Pop()
			currentBlock.appendChild(innerBlock)
		case ';':
			// End the current word and statement.
			currentStmt = append(currentStmt, string(currentWord))
			currentWord = nil

			if len(currentStmt) > 0 {
				if node, err = newNode(currentStmt); err != nil {
					return nil, err
				}
				currentBlock.appendChild(node)
			}
			currentStmt = nil
		case '\n', ' ', '\t', '\r':
			// End the current word.
			if len(currentWord) > 0 {
				currentStmt = append(currentStmt, string(currentWord))
				currentWord = nil
			}
		default:
			// Add current character onto the current word.
			currentWord = append(currentWord, char)
		}
	}

	return currentBlock, nil
}

func DecodeFile(filePath string) (*Node, error) {
	fileContent, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(fileContent) <= 0 {
		return nil, errors.New("file content is empty")
	}
	return Decode(fileContent)
}

// Read file
func readFile(path string) (string, error) {
	contentByte, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contentByte), nil
}

// Judge whether the file exists, return true if it exists, and return false if it does not exist
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
