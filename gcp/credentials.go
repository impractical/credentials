package gcp

import (
	"context"
	"fmt"

	"impractical.co/credentials"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	_ credentials.Accessor = &Credentials{}
	_ credentials.Setter   = &Credentials{}
	_ credentials.Initer   = &Credentials{}
	_ credentials.Closer   = &Credentials{}
)

type Credentials struct {
	Project string
	Version string
	// TODO: auth options

	client *secretmanager.Client
}

func (c *Credentials) newSecretClient(ctx context.Context) (*secretmanager.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	// TODO: support other auth options
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}

func (c *Credentials) Init(ctx context.Context) error {
	_, err := c.newSecretClient(ctx)
	return err
}

func (c *Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	client, err := c.newSecretClient(ctx)
	if err != nil {
		return nil, err
	}
	version := "latest"
	if c.Version != "" {
		version = c.Version
	}
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", c.Project, id, version)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.Payload.Data, nil
}

func (c *Credentials) Set(ctx context.Context, id string, plaintext []byte) error {
	client, err := c.newSecretClient(ctx)
	if err != nil {
		return err
	}

	parent := fmt.Sprintf("projects/%s/secrets/%s", c.Project, id)
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: plaintext,
		},
	}

	_, err = client.AddSecretVersion(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Credentials) Close(_ context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}
