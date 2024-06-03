package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
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

	timeNow := time.Now().Format(time.RFC3339)

	if err := s.driver.QueryRowContext(newCtx, createPost,
		username, post.Name, post.Description, post.Content,
		post.CommentsAllowed, timeNow, timeNow).
		Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	const op = "internal.storage.postgres.getPosts"

	newCtx, cancel := context.WithTimeout(ctx, timeToAnswer)
	defer cancel()

	rows, err := s.driver.QueryContext(newCtx, getPosts, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var posts = make([]*model.Post, 0, limit)
	var ch = make(chan *sql.Rows)
	var doneCh = make(chan struct{})
	go func() {
		for {
			select {
			case rows, opened := <-ch:
				if !opened {
					return
				}
				var post model.Post
				err = rows.Scan(&post.ID, &post.Author, &post.Name, &post.Description, &post.Content, &post.CommentsAllowed, &post.CreatedAt, &post.UpdatedAt)
				if err != nil {
					slog.Error("couldn't scan post for posts request", slogAttr.SlogErr(op, err))
					doneCh <- struct{}{}
					continue
				}
				posts = append(posts, &post)
				doneCh <- struct{}{}
			}
		}
	}()

	for rows.Next() {
		ch <- rows
		<-doneCh
	}
	close(ch)
	close(doneCh)
	return posts, nil
}
