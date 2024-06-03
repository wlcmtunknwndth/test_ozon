package inmemory

import (
	"github.com/arriqaaq/flashdb"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"sync"
)

type comments struct {
	storage map[string]model.Comment
	mtx     sync.RWMutex
}

type posts struct {
	storage map[string]model.Post
	mtx     sync.RWMutex
}

type users struct {
	storage map[string]model.User
	mtx     sync.RWMutex
}

type Storage struct {
	comments comments
	posts    posts
	users    users
}

func New(config *flashdb.Config) *Storage {
	const op = "internal.storage.inmemory.New"
	return &Storage{}
}
