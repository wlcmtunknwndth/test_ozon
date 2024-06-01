package postgres

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"sync"
	"time"
)

const (
	timeToAnswer = 15 * time.Second
)

func (s *Storage) CreatePost(ctx context.Context, username string, post *model.NewPost) (uint64, error) {
	const op = "internal.storage.postgres.CreatePost"
	newCtx, cancel := context.WithTimeout(ctx, timeToAnswer)
	defer cancel()
	var id uint64

	if err := s.driver.QueryRowContext(newCtx, createPost,
		username, post.Name, post.Description, post.Content,
		post.CommentsAllowed, time.Now().Format(time.RFC3339)).
		Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) getPosts(ctx context.Context, limit, offset int) ([]model.Post, error) {
	const op = "internal.storage.postgres.getPosts"

	newCtx, cancel := context.WithTimeout(ctx, timeToAnswer)
	defer cancel()

	rows, err := s.driver.QueryContext(newCtx, getPosts, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var posts = make([]model.Post, 0, limit)
	var wg sync.WaitGroup
	for rows.Next(){
		wg.Add(1)
		go func() {
			defer wg.Done()
			var post model.Post
			rows.Scan(&post.ID, &post.)

		}()
		wg.Wait()
	}
}

func (s *Storage) CreateComment(ctx context.Context, username string, comment *model.NewComment) (uint64, error) {
	const op = "internal.storage.postgres.CreateComment"
	newCtx, cancel := context.WithTimeout(ctx, timeToAnswer)
	defer cancel()

	var id uint64
	if err := s.driver.QueryRowContext(newCtx, createComment, username, comment.PostID,
		comment.RepliesTo, comment.Text,
		time.Now().Format(time.RFC3339)).
		Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
