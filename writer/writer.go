package writer

type Writer interface {
	Write(data []byte) error

	//程序结束时，注意close
	Close() error
}

type Builder interface {
	Build() (Writer, error)
}
