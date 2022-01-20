package offset_record

type OffsetRecord struct {
	FilePath string `json:"filePath"`
	Offset   int64  `json:"offset"`
}
