package writer

type Writer interface {
	Write(data []byte) error
}

type Builder interface {
	Build() (Writer, error)
}
