package collector

import (
	"context"
	"log"
	"log-collector/reader"
	"log-collector/writer"
)

type Collector struct {
	reader  []reader.Reader
	writer  []writer.Writer
	MsgChan chan []byte
}

func NewCollector(reader []reader.Reader, writer []writer.Writer, num uint) *Collector {
	return &Collector{
		reader:  reader,
		writer:  writer,
		MsgChan: make(chan []byte, num),
	}
}
func (c *Collector) Collect(ctx context.Context) error {
	errch := make(chan error)
	//注意结束时把writer给关闭
	defer func() {
		close(errch)
		for _, w := range c.writer {
			w.Close()
		}
	}()

	//reader返回错误都是大型错误,需要监听
	for _, r := range c.reader {
		go func(chan error) {
			errch <- r.Read(ctx, c.MsgChan)
		}(errch)
	}

	go c.write(ctx)

	//从errch中获取错误,如果有错误就返回，告知主程序取消
	for {
		select {
		case err := <-errch:
			return err
		}
	}
}
func (c *Collector) write(ctx context.Context) {
	for {
		select {
		case msg := <-c.MsgChan:
			for _, w := range c.writer {
				go func() {
					err := w.Write(msg)
					if err != nil {
						log.Println(err)
					}
				}()
			}
		case <-ctx.Done():
			return
		}
	}
}
