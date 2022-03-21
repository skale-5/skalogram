package s3

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/skale-5/skalogram/web"
)

type Client struct {
	sess     *session.Session
	s3       *s3.S3
	uploader *s3manager.Uploader
}

func NewClient(ctx context.Context, region string) *Client {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: &region,
			},
		},
	)
	if err != nil {
		log.Fatal("cannot initialize s3 session:", err)
	}

	return &Client{
		sess:     sess,
		s3:       s3.New(sess),
		uploader: s3manager.NewUploader(sess),
	}
}

func (c *Client) Write(ctx context.Context, object *web.ObjectPath, content io.Reader) error {
	_, err := c.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Path),
		Body:   content,
	})

	if err != nil {
		return fmt.Errorf("failed to upload s3 %s: %s", object.URL(), err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, object *web.ObjectPath) (io.ReadCloser, error) {
	out, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get s3 object %s: %s", object.URL(), err)
	}

	return out.Body, nil
}
