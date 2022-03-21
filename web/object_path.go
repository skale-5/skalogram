package web

import (
	"fmt"
	"net/url"
	"strings"
)

type ObjectPath struct {
	Scheme string
	Bucket string
	Path   string
}

func (op *ObjectPath) URL() string {
	return fmt.Sprintf("%s://%s/%s", op.Scheme, op.Bucket, op.Path)
}

func NewObjectPath(objectURL string) (*ObjectPath, error) {
	u, err := url.Parse(objectURL)
	if err != nil {
		return nil, err
	}
	op := &ObjectPath{}
	op.Scheme = u.Scheme
	op.Bucket = u.Host
	op.Path = strings.TrimPrefix(u.Path, "/")
	return op, nil
}
