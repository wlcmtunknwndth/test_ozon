package graph

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate
type Storage interface {
	CreateComment(ctx context.Context, username string, comment *model.NewComment) (uint64, error)
	GetComments(ctx context.Context, postId uint64) ([]*model.Comment, error)
	CreatePost(ctx context.Context, username string, post *model.NewPost) (uint64, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error)
	GetPassword(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, usr *model.NewUser) error
	IsAdmin(ctx context.Context, username string) (bool, error)
}

type Resolver struct {
	Storage Storage
}
