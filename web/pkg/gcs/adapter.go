package gcs

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
	"github.com/skale-5/skalogram/web"
)

type Client struct {
	storage *storage.Client
}

func NewClient(ctx context.Context) *Client {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("cannot create GCS client:", err)
	}

	return &Client{
		storage: client,
	}
}

func (c *Client) Write(ctx context.Context, object *web.ObjectPath, content io.Reader) error {
	obj := c.storage.Bucket(object.Bucket).Object(object.Path)
	w := obj.NewWriter(ctx)
	defer w.Close()

	_, err := io.Copy(w, content)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, object *web.ObjectPath) (io.ReadCloser, error) {
	obj := c.storage.Bucket(object.Bucket).Object(object.Path)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}
