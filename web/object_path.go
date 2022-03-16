package web

import (
	"fmt"
	"net/url"
	"strings"
)

type ObjectPath struct {
	Bucket string
	Path   string
}

func (op *ObjectPath) URL() string {
	return fmt.Sprintf("gs://%s/%s", op.Bucket, op.Path)
}

func NewObjectPath(objectURL string) (*ObjectPath, error) {
	u, err := url.Parse(objectURL)
	if err != nil {
		return nil, err
	}
	op := &ObjectPath{}
	op.Bucket = u.Host
	op.Path = strings.TrimPrefix(u.Path, "/")
	return op, nil
}
