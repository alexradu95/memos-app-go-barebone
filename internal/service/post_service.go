package service

import (
	"context"
	"journal-lite/internal/posts"
	"journal-lite/internal/repository"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, post posts.Post) (posts.Post, error) {
	return s.repo.CreatePost(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, postId int64) error {
	return s.repo.DeletePost(ctx, postId)
}

func (s *PostService) GetPosts(ctx context.Context, params posts.QueryParams) ([]posts.Post, error) {
	return s.repo.GetPosts(ctx, params)
}

func (s *PostService) GetPost(ctx context.Context, userId int64, postId int64) (posts.Post, error) {
	return s.repo.GetPost(ctx, userId, postId)
}

func (s *PostService) UpdatePost(ctx context.Context, newContent string, postId int64) error {
	return s.repo.UpdatePost(ctx, newContent, postId)
}
