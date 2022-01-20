package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var AlreadyExists = "The folder or file already exists"

// RemoveDir Delete directory
func RemoveDir(path string) {
	// Delete file
	dir := path
	exist, err := PathExists(dir)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if exist {
			// os.RemoveAll is traversal deletion, and folders and files can be used
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Println(dir+"Deletion failed：", err.Error())
			} else {
				fmt.Println(dir + "Delete succeeded！")
			}
		} else {
			fmt.Println(dir + "File, folder does not exist！")
		}
	}
}

// RemoveFile Delete file
func RemoveFile(filePath string) error {
	var err error
	if IsExists(filePath, true) {
		err = os.Remove(filePath)
	}
	return err
}

// WriteFile Write file
func WriteFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), 0777)
	if err != nil {
		return err
	}
	return nil
}

// WriteFileByte Write file
func WriteFileByte(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, 0777)
	if err != nil {
		return err
	}
	return nil
}

// WriteFileAdv When a file is written, a folder is created
func WriteFileAdv(path string, content string) error {
	dir := filepath.Dir(path)
	if !IsExists(dir, false) {
		err := CreateFolder(dir)
		if err != nil {
			return err
		}
	}
	return WriteFile(path, content)
}

// PathExists Determine whether the folder exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsExists Determine whether the file or folder exists
func IsExists(path string, isFile bool) bool {
	s, err := os.Stat(path)
	if err == nil {
		if s.IsDir() == !isFile {
			return true
		}
	}
	return false
}

// CreateFolder Create folder
func CreateFolder(filePath string) error {
	if IsExists(filePath, false) {
		return errors.New(AlreadyExists)
	}
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(filePath, 0777)
	if err != nil {
		return err
	}
	return nil
}

// ReadFile Read file
func ReadFile(path string) (string, error) {
	contentByte, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contentByte), nil
}

// ReadFileAndToInt64 Read the contents of the file and convert it to int64
func ReadFileAndToInt64(path string) (int64, error) {
	str, err := ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(str) <= 0 {
		return 0, errors.New("file content is empty")
	}
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

// ReadFileByte Read file and return byte array
func ReadFileByte(path string) ([]byte, error) {
	contentByte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return contentByte, nil
}

// GetFileBaseName Get the basic name of the file and remove the specified suffix
func GetFileBaseName(filePath string, trimSuffix string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, trimSuffix)
}
