// internal/repository/sqlite/post_repository.go
package sqlite

import (
	"context"
	"database/sql"
	"journal-lite/internal/posts"
	"journal-lite/internal/repository"
	"time"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, post posts.Post) (posts.Post, error) {
	query := `INSERT INTO posts (content, created_at, updated_at, account_id) VALUES (?, ?, ?, ?) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, post.Content, post.CreatedAt, post.UpdatedAt, post.AccountId).Scan(&post.Id)
	return post, err
}

func (r *PostRepository) DeletePost(ctx context.Context, postId int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", postId)
	return err
}

func (r *PostRepository) GetPosts(ctx context.Context, params posts.QueryParams) ([]posts.Post, error) {
	args := []interface{}{params.AccountId}
	query := `SELECT id, content, created_at, updated_at, account_id FROM posts WHERE account_id = ?`

	if params.SearchText != "" {
		query += " AND content LIKE ?"
		args = append(args, "%"+params.SearchText+"%")
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err != nil {
			return nil, err
		}
		query += " AND created_at >= ?"
		args = append(args, dateFrom)
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err != nil {
			return nil, err
		}
		query += " AND created_at <= ?"
		args = append(args, dateTo)
	}

	query += " ORDER BY created_at DESC"

	if params.PageNumber > 0 && params.PageSize > 0 {
		offset := (params.PageNumber - 1) * params.PageSize
		query += " LIMIT ? OFFSET ?"
		args = append(args, params.PageSize, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []posts.Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Id, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.AccountId); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetPost(ctx context.Context, userId int64, postId int64) (posts.Post, error) {
	var post posts.Post
	err := r.db.QueryRowContext(ctx, "SELECT id, content, created_at, updated_at, account_id FROM posts WHERE id = ? AND account_id = ?", postId, userId).
		Scan(&post.Id, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.AccountId)
	return post, err
}

func (r *PostRepository) UpdatePost(ctx context.Context, newContent string, postId int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE posts SET content = ?, updated_at = ? WHERE id = ?", newContent, time.Now(), postId)
	return err
}
