package writer

import (
	"testing"
)

func TestStdoutWriter_Write(t *testing.T) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{[]byte("this is a test")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StdoutWriter{}
			if err := s.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
