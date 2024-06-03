package inmemory

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"slices"
	"time"
)

func (s *Storage) CreateComment(ctx context.Context, username string, comment *model.NewComment) (uint64, error) {
	const op = "internal.storage.inmemory.comment.GetComments"
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("%s: time limie exceeded", op)
		default:
			id := s.comments.counter.Add(1)
			s.posts.storage.Store(time.Now().UnixNano(), model.Comment{
				ID:        id,
				PostID:    comment.PostID,
				RepliesTo: comment.RepliesTo,
				Author:    username,
				Text:      comment.Text,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			return id, nil
		}
	}
}

func (s *Storage) GetComments(ctx context.Context, postId uint64) ([]*model.Comment, error) {
	const op = "internal.storage.inmemory.comment.GetComments"
	arr := make([]*model.Comment, 0, 4)
	var i = 0
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%s: time limit exceeded", op)
		default:
			s.posts.storage.Range(func(key, value any) bool {
				defer func() {
					i = i + 1
				}()
				comment, ok := value.(model.Comment)
				if !ok {
					return true
				}
				if len(arr) == 0 {
					arr = append(arr, &comment)
					return true
				}
				if comment.CreatedAt.UnixNano() > arr[len(arr)-1].CreatedAt.UnixNano() {
					arr = append(arr, arr[len(arr)-1])
					arr[len(arr)-2] = &comment
				} else {
					arr = append(arr, &comment)
				}
				return true
			})
			slices.SortFunc(arr, func(a, b *model.Comment) int {
				if a.CreatedAt.UnixNano() > b.CreatedAt.UnixNano() {
					return 1
				} else if a.CreatedAt.UnixNano() < b.CreatedAt.UnixNano() {
					return -1
				} else {
					return 0
				}
			})
			return arr, nil
		}
	}
}
