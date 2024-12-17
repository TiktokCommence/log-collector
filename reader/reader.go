package reader

import "context"

type Reader interface {
	Read(ctx context.Context, ch chan<- []byte) error
}

type Builder interface {
	Build() (Reader, error)
}
