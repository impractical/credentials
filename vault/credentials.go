package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"

	"impractical.co/credentials"
)

var (
	_ credentials.Accessor = &Credentials{}
	_ credentials.Initer   = &Credentials{}
)

type Credentials struct {
	Address   string
	MountPath string
	Object    string
	Token     string

	client *vault.Client
}

func (c Credentials) newSecretClient(ctx context.Context) (*vault.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	client, err := vault.NewClient(&vault.Config{
		Address: c.Address,
	})
	if err != nil {
		return nil, err
	}
	client.SetToken(c.Token)
	c.client = client
	return client, nil
}

func (c Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	client, err := c.newSecretClient(ctx)
	if err != nil {
		return nil, err
	}
	key := c.MountPath + "/data/" + c.Object
	data, err := client.Logical().Read(key)
	if err != nil {
		return nil, err
	}
	dataI, ok := data.Data["data"]
	if !ok {
		return nil, fmt.Errorf("No data key in response from Vault for %s", key)
	}
	dataMap, ok := dataI.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("Data key in response from Vault for %s wasn't map[string]any, it was %T", key, dataI)
	}
	iface, ok := dataMap[id]
	if !ok {
		return nil, fmt.Errorf("No %s entry found in Vault for %s", id, key)
	}
	val, ok := iface.(string)
	if !ok {
		return nil, fmt.Errorf("%s data set in %s in Vault wasn't string, it was %T", id, key, iface)
	}
	return []byte(val), nil
}

func (c *Credentials) Init(ctx context.Context) error {
	_, err := c.newSecretClient(ctx)
	return err
}
