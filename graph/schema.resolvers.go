package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"

	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"github.com/wlcmtunknwndth/test_ozon/internal/auth"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input *model.NewPost) (uint64, error) {
	const op = "graph.schema.resolvers.CreatePost"
	if !auth.IsRegistered(ctx) {
		return 0, fmt.Errorf("%s: %s", op, notAuthenticated)
	}
	username := auth.GetUsername(ctx)
	if username == "" {
		return 0, fmt.Errorf("%s: %s", op, noUsername)
	}
	id, err := r.Storage.CreatePost(ctx, username, input)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input *model.NewComment) (uint64, error) {
	const op = "graph.schema.resolvers.CreateComment"
	if !auth.IsRegistered(ctx) {
		return 0, fmt.Errorf("%s: %s", op, notAuthenticated)
	}

	username := auth.GetUsername(ctx)
	if username == "" {
		return 0, fmt.Errorf("%s: %s", op, noUsername)
	}

	id, err := r.Storage.CreateComment(ctx, username, input)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Posts is the resolver for the Posts field.
func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	const op = "graph.schema.resolvers.Posts"
	//if !auth.IsRegistered(ctx){
	//	return 0, fmt.Errorf("%s: %s", op, notAuthenticated)
	//}

	posts, err := r.Storage.GetPosts(ctx, *limit, *offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

// Comments is the resolver for the Comments field.
func (r *queryResolver) Comments(ctx context.Context, postID *uint64) ([]*model.Comment, error) {
	const op = "graph.schema.resolvers.Comments"

	posts, err := r.Storage.GetComments(ctx, *postID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
const (
	notAuthenticated = "not authenticated"
	noUsername       = "couldn't get username"
)
