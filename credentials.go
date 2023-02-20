package credentials

import "context"

type Accessor interface {
	Get(context.Context, string) ([]byte, error)
}

type Setter interface {
	Set(context.Context, string, []byte) error
}

type Initer interface {
	Init(context.Context) error
}

type Closer interface {
	Close(context.Context) error
}
