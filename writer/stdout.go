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
	fmt.Println(string(data))
	return nil
}