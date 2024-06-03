package inmemory

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"slices"
	"time"
)

func (s *Storage) CreatePost(ctx context.Context, username string, post *model.NewPost) (uint64, error) {
	const op = "internal.storage.inmemory.posts"
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("%s: time limie exceeded", op)
		default:
			id := s.posts.counter.Add(1)
			s.posts.storage.Store(time.Now().UnixNano(), model.Post{
				ID:              id,
				Author:          username,
				Name:            post.Name,
				Description:     post.Description,
				Content:         post.Content,
				CommentsAllowed: post.CommentsAllowed,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			})

			return id, nil
		}
	}
}

func (s *Storage) GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	arr := make([]*model.Post, 0, limit)
	var i = 0
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return arr, nil
		default:
			s.posts.storage.Range(func(key, value any) bool {
				defer func() {
					i = i + 1
				}()
				if i < offset {
					return true
				}
				if i >= limit+offset {
					return false
				}
				post, ok := value.(model.Post)
				if !ok {
					return true
				}
				if len(arr) == 0 {
					arr = append(arr, &post)
					return true
				}
				if post.CreatedAt.UnixNano() > arr[len(arr)-1].CreatedAt.UnixNano() {
					arr = append(arr, arr[len(arr)-1])
					arr[len(arr)-2] = &post
				} else {
					arr = append(arr, &post)
				}
				return true
			})
			slices.SortFunc(arr, func(a, b *model.Post) int {
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
