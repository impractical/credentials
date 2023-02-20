package envvar

import (
	"context"
	"fmt"
	"os"

	"impractical.co/credentials"
)

var (
	_ credentials.Accessor = Credentials{}
)

type Credentials struct {
	Prefix string
}

func (c Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	if v := os.Getenv(c.Prefix + id); v != "" {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("env var %q not set", c.Prefix+id)
}
