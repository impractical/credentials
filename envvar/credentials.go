package envvar

import (
	"context"
	"fmt"
	"os"
	"strings"

	"impractical.co/credentials"
)

var (
	_ credentials.Accessor = Credentials{}
)

type Credentials struct {
	Prefix      string
	CoerceUpper bool
}

func (c Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	env := c.Prefix + id
	if c.CoerceUpper {
		env = strings.ToUpper(env)
	}
	if v := os.Getenv(c.Prefix + id); v != "" {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("env var %q not set", c.Prefix+id)
}
