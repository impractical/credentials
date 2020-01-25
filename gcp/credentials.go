package gcp

import (
	"context"
	"fmt"

	"impractical.co/credentials"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

var (
	_ credentials.Accessor = Credentials{}
	_ credentials.Setter   = Credentials{}
)

type Credentials struct {
	Project string
	Version string
	// TODO: auth options
}

func (c Credentials) newSecretClient(ctx context.Context) (*secretmanager.Client, error) {
	// TODO: support other auth options
	return secretmanager.NewClient(ctx)
}

func (c Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	client, err := c.newSecretClient(ctx)
	if err != nil {
		// TODO: better error handling
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
		// TODO: better error handling
		return nil, err
	}
	return result.Payload.Data, nil
}

func (c Credentials) Set(ctx context.Context, id string, plaintext []byte) error {
	client, err := c.newSecretClient(ctx)
	if err != nil {
		// TODO: better error handling
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
		// TODO: better error handling
		return err
	}
	return nil
}
