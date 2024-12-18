package writer

import (
	"fmt"
	"sync"
)

//将日志写入到控制台

type StdoutWriter struct {
	mutex sync.Mutex
}

func (s *StdoutWriter) Write(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fmt.Printf(string(data))
	return nil
}
func (s *StdoutWriter) Close() error {
	return nil
}
