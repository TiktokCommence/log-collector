package writer

type FileWriterBuilder struct {
	FilePath     string
	FileName     string
	MaxSize      int64
	RotateByTime bool //是否根据时间来进行切割
}

func NewFileWriterBuilder(filePath string, fileName string, maxSize int64, rotateByTime bool) *FileWriterBuilder {
	return &FileWriterBuilder{
		FilePath:     filePath,
		FileName:     fileName,
		MaxSize:      maxSize,
		RotateByTime: rotateByTime,
	}
}
func (f *FileWriterBuilder) Build() (Writer, error) {
	return &FileWriter{
		filePath:     f.FilePath,
		filename:     f.FileName,
		maxSize:      f.MaxSize,
		rotateByTime: f.RotateByTime,
	}, nil
}
