package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"sync"
	"time"
)

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

func (s *Storage) GetComments(ctx context.Context, postId uint64) ([]model.Comment, error) {
	const op = "internal.storage.postgres.GetComments"

	newCtx, cancel := context.WithTimeout(ctx, timeToAnswer)
	defer cancel()

	rows, err := s.driver.QueryContext(newCtx, getComments, postId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var comments []model.Comment
	var ch = make(chan *sql.Rows)
	var mtx sync.Mutex
	mtx.Unlock()
	go func() {
		defer mtx.Unlock()
		for {
			select {
			case row, opened := <-ch:
				if !opened {
					return
				}
				var comment model.Comment
				if err := row.Scan(&comment.ID, &comment.Author,
					&comment.PostID, &comment.RepliesTo, &comment.Text,
					&comment.CreatedAt, &comment.UpdatedAt); err != nil {
					slog.Error("couldn't scan comment", slogAttr.SlogErr(op, err))
					mtx.Unlock()
					continue
				}
				comments = append(comments, comment)
				mtx.Unlock()
			}
		}
	}()
	for rows.Next() {
		ch <- rows
		mtx.Lock()
	}
	close(ch)
	return comments, nil
}
