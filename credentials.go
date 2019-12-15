package credentials

import "context"

type Accessor interface {
	Get(context.Context, string) ([]byte, error)
}

type Setter interface {
	Set(context.Context, string, []byte) error
}
