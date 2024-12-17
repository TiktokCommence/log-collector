package writer

type StdoutWriterBuilder struct{}

func (s *StdoutWriterBuilder) Build() (Writer, error) {
	return &StdoutWriter{}, nil
}
