package gcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"impractical.co/credentials"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/storage"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var (
	_ credentials.Accessor = Credentials{}
	_ credentials.Setter   = Credentials{}
)

type Credentials struct {
	StorageBucket string
	KMSKey        struct {
		Project string
		RingID  string
		KeyID   string
	}

	// TODO: auth options
}

func (c Credentials) newStorageClient(ctx context.Context) (*storage.Client, error) {
	// TODO: support other auth options
	return storage.NewClient(ctx)
}

func (c Credentials) newKMSClient(ctx context.Context) (*kms.KeyManagementClient, error) {
	// TODO: support other auth options
	return kms.NewKeyManagementClient(ctx)
}

func (c Credentials) Get(ctx context.Context, id string) ([]byte, error) {
	storageClient, err := c.newStorageClient(ctx)
	if err != nil {
		// TODO: better error handling
		return nil, err
	}
	bucket := storageClient.Bucket(c.StorageBucket)
	object := bucket.Object(id)
	rc, err := object.NewReader(ctx)
	if err != nil {
		// TODO: better error handling
		return nil, err
	}
	cipherBytes, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		// TODO: better error handling
		return nil, err
	}

	cipherb64 := strings.TrimSpace(string(cipherBytes))

	ciphertext, err := base64.StdEncoding.DecodeString(cipherb64)
	if err != nil {
		// TODO: better error handling
		return nil, err
	}
	kmsClient, err := c.newKMSClient(ctx)
	if err != nil {
		// TODO: better error handling
		return nil, err
	}

	req := &kmspb.DecryptRequest{
		Name:       fmt.Sprintf("projects/%s/locations/global/keyRings/%s/cryptoKeys/%s", c.KMSKey.Project, c.KMSKey.RingID, c.KMSKey.KeyID),
		Ciphertext: ciphertext,
	}

	resp, err := kmsClient.Decrypt(ctx, req)
	if err != nil {
		// TODO: better error handling
		return nil, err
	}
	return resp.Plaintext, nil
}

func (c Credentials) Set(ctx context.Context, id string, plaintext []byte) error {
	kmsClient, err := c.newKMSClient(ctx)
	if err != nil {
		// TODO: better error handling
		return err
	}
	storageClient, err := c.newStorageClient(ctx)
	if err != nil {
		// TODO: better error handling
		return err
	}

	req := &kmspb.EncryptRequest{
		Name:      fmt.Sprintf("projects/%s/locations/global/keyRings/%s/cryptoKeys/%s", c.KMSKey.Project, c.KMSKey.RingID, c.KMSKey.KeyID),
		Plaintext: plaintext,
	}
	resp, err := kmsClient.Encrypt(ctx, req)
	if err != nil {
		// TODO: better error handling
		return err
	}
	cipherb64 := base64.StdEncoding.EncodeToString(resp.Ciphertext)

	bucket := storageClient.Bucket(c.StorageBucket)
	object := bucket.Object(id)
	w := object.NewWriter(ctx)

	_, err = fmt.Fprintf(w, cipherb64)
	if err != nil {
		// TODO: better error handling
		return err
	}
	err = w.Close()
	if err != nil {
		// TODO: better error handling
		return err
	}
	return nil
}
