package envvar

import (
	"context"
	"fmt"
	"os"

	"impractical.co/credentials"
)

var (
	_ credentials.Accessor = Credentials{}
	_ credentials.Setter   = Credentials{}
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

func (c Credentials) Set(ctx context.Context, id string, plaintext []byte) error {
	return fmt.Errorf("setting credentials is not supported")
}
