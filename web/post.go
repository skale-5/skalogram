package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/qeesung/image2ascii/convert"
	"github.com/robert-nix/ansihtml"
)

type Post struct {
	ID        uuid.UUID
	Score     int
	ImgUrl    string
	CreatedAt time.Time
}

func GenerateAscii(file io.ReadCloser) (string, error) {
	img, _, err := image.Decode(file)
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
	ID     uuid.UUID
	ImgUrl string
}

type PostDatabaseAdapter interface {
	CreateTable(ctx context.Context) error
	CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	GetPost(ctx context.Context, id uuid.UUID) (Post, error)
	ListPosts(ctx context.Context) ([]Post, error)
	UpvotePost(ctx context.Context, id uuid.UUID) error
	DownvotePost(ctx context.Context, id uuid.UUID) error
}

var ErrPostCacheNotFound = errors.New("post not found in cache")

type PostCacheAdapter interface {
	Ping(ctx context.Context) error
	CachePost(ctx context.Context, id uuid.UUID, content interface{}, ttl time.Duration) (interface{}, error)
	GetPost(ctx context.Context, id uuid.UUID) (interface{}, error)
}

type PostStorageAdapter interface {
	Write(ctx context.Context, object *ObjectPath, content io.Reader) error
	Get(ctx context.Context, object *ObjectPath) (io.ReadCloser, error)
}

type PostCacheService struct {
	adapter PostCacheAdapter
}

func NewPostCacheService(a PostCacheAdapter) *PostCacheService {
	return &PostCacheService{
		adapter: a,
	}
}

func (pcs *PostCacheService) Ping(ctx context.Context) error {
	return pcs.adapter.Ping(ctx)
}

func (pcs *PostCacheService) CachePost(ctx context.Context, id uuid.UUID, content string, ttl time.Duration) (string, error) {
	ret, err := pcs.adapter.CachePost(ctx, id, content, ttl)
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
	ret, err := pcs.adapter.GetPost(ctx, id)
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
	adapter PostDatabaseAdapter
}

func NewPostDatabaseService(a PostDatabaseAdapter) *PostDatabaseService {
	return &PostDatabaseService{
		adapter: a,
	}
}

func (pds *PostDatabaseService) CreateTable(ctx context.Context) error {
	err := pds.adapter.CreateTable(ctx)
	if err != nil {
		return fmt.Errorf("cannot create post table: %w", err)
	}
	return nil
}

func (pds *PostDatabaseService) ListPosts(ctx context.Context) ([]Post, error) {
	posts, err := pds.adapter.ListPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot list posts: %w", err)
	}
	return posts, nil
}

func (pds *PostDatabaseService) UpvotePost(ctx context.Context, id uuid.UUID) error {
	err := pds.adapter.UpvotePost(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot upvote post: %w", err)
	}
	return nil
}

func (pds *PostDatabaseService) DownvotePost(ctx context.Context, id uuid.UUID) error {
	err := pds.adapter.DownvotePost(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot downvote post: %w", err)
	}
	return nil
}

func (pds *PostDatabaseService) CreatePost(ctx context.Context, args CreatePostParams) error {
	_, err := pds.adapter.CreatePost(ctx, args)
	if err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}
	return nil
}

type PostStorageService struct {
	adapter PostStorageAdapter
}

func NewPostStorageService(a PostStorageAdapter) *PostStorageService {
	return &PostStorageService{
		adapter: a,
	}
}

func (pss *PostStorageService) Write(ctx context.Context, object *ObjectPath, r io.Reader) error {
	return pss.adapter.Write(ctx, object, r)
}

func (pss *PostStorageService) Get(ctx context.Context, object *ObjectPath) (io.ReadCloser, error) {
	return pss.adapter.Get(ctx, object)
}
