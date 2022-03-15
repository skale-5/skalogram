package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/qeesung/image2ascii/convert"
	"github.com/robert-nix/ansihtml"
)

/*
CREATE TABLE posts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	score INTEGER NOT NULL DEFAULT 0,
	img_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
*/

type Post struct {
	ID        uuid.UUID
	Score     int
	ImgUrl    string
	CreatedAt time.Time
}

func (p *Post) GenerateAscii() (string, error) {

	resp, err := http.Get(p.ImgUrl)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	converter := convert.NewImageConverter()
	ascii := converter.Image2ASCIIString(img, &convert.Options{
		FixedWidth:  55,
		FixedHeight: 20,
		Ratio:       2,
		Colored:     true,
	})

	html := ansihtml.ConvertToHTML([]byte(ascii))
	return string(html), nil
}

type CreatePostParams struct {
	Score  string
	ImgUrl string
}

type PostDatabaseAdapter interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	GetPost(ctx context.Context, id uuid.UUID) (Post, error)
	ListPosts(ctx context.Context) ([]Post, error)
	UpvotePost(ctx context.Context, id uuid.UUID) error
	DownvotePost(ctx context.Context, id uuid.UUID) error
}

var ErrPostCacheNotFound = errors.New("post not found in cache")

type PostCacheAdapter interface {
	CachePost(ctx context.Context, id uuid.UUID, content interface{}, ttl time.Duration) (interface{}, error)
	GetPost(ctx context.Context, id uuid.UUID) (interface{}, error)
}

type PostStorageAdapter interface {
	Write(ctx context.Context, id uuid.UUID, content io.Reader) (string, error)
	Get(ctx context.Context, id uuid.UUID) (io.Reader, error)
}

type PostCacheService struct {
	Adapter PostCacheAdapter
}

func (pcs *PostCacheService) CachePost(ctx context.Context, id uuid.UUID, content string, ttl time.Duration) (string, error) {
	ret, err := pcs.Adapter.CachePost(ctx, id, content, ttl)
	if err != nil {
		return "", err
	}
	content, ok := ret.(string)
	if !ok {
		return "", ErrPostCacheNotFound
	}
	return content, nil
}

func (pcs *PostCacheService) GetPost(ctx context.Context, id uuid.UUID) (string, error) {
	ret, err := pcs.Adapter.GetPost(ctx, id)
	if err != nil {
		return "", err
	}
	content, ok := ret.(string)
	if !ok {
		return "", ErrPostCacheNotFound
	}
	return content, nil
}

type PostDatabaseService struct {
	Adapter PostDatabaseAdapter
}

func (pds *PostDatabaseService) ListPosts(ctx context.Context) ([]Post, error) {
	posts, err := pds.Adapter.ListPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot list posts: %w", err)
	}
	return posts, nil
}

func (pds *PostDatabaseService) UpvotePost(ctx context.Context, id uuid.UUID) error {
	err := pds.Adapter.UpvotePost(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot upvote post: %w", err)
	}
	return nil
}

func (pds *PostDatabaseService) DownvotePost(ctx context.Context, id uuid.UUID) error {
	err := pds.Adapter.DownvotePost(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot downvote post: %w", err)
	}
	return nil
}
