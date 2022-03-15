package post

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/skale-5/skalogram/web"
)

const createPost = `-- name: CreatePost :execresult
INSERT INTO posts (
  score, img_url
) VALUES (
  $1, $1
)
`

func (q *Queries) CreatePost(ctx context.Context, arg web.CreatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createPost, arg.Score, arg.ImgUrl)
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPost = `-- name: GetPost :one
SELECT id, score, img_url, created_at FROM posts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id uuid.UUID) (web.Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i web.Post
	err := row.Scan(
		&i.ID,
		&i.Score,
		&i.ImgUrl,
		&i.CreatedAt,
	)
	return i, err
}

const listPosts = `-- name: ListPosts :many
SELECT id, score, img_url, created_at FROM posts
ORDER BY created_at DESC
`

func (q *Queries) ListPosts(ctx context.Context) ([]web.Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []web.Post
	for rows.Next() {
		var i web.Post
		if err := rows.Scan(
			&i.ID,
			&i.Score,
			&i.ImgUrl,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upvotePost = `-- name: UpvotePost :exec
UPDATE posts
SET score = score + 1
WHERE id = $1
`

func (q *Queries) UpvotePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upvotePost, id)
	return err
}

const downvotePost = `-- name: DownvotePost :exec
UPDATE posts SET score = score - 1
WHERE id = $1
`

func (q *Queries) DownvotePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, downvotePost, id)
	return err
}
