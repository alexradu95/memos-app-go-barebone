package repository

import (
	"context"
	"journal-lite/internal/posts"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post posts.Post) (posts.Post, error)
	DeletePost(ctx context.Context, postId int64) error
	GetPosts(ctx context.Context, params posts.QueryParams) ([]posts.Post, error)
	GetPost(ctx context.Context, userId int64, postId int64) (posts.Post, error)
	UpdatePost(ctx context.Context, newContent string, postId int64) error
}
