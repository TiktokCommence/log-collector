package reader

import (
	"context"
	"testing"
)

func TestRead(t *testing.T) {
	builder := NewKafkaReaderBuilder([]string{"127.0.0.1:9092"}, "testlog")
	reader, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	var MsgCh = make(chan []byte, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx2 context.Context) {
		for {
			select {
			case <-ctx2.Done():
				return
			case msg := <-MsgCh:
				t.Log(string(msg))
			}
		}
	}(ctx)

	err = reader.Read(ctx, MsgCh)
	if err != nil {
		t.Fatal(err)
	}
}
