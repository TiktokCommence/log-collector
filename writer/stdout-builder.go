package writer

type StdoutWriterBuilder struct{}

func NewStdoutWriterBuilder() *StdoutWriterBuilder {
	return &StdoutWriterBuilder{}
}

func (s *StdoutWriterBuilder) Build() (Writer, error) {
	return &StdoutWriter{}, nil
}
