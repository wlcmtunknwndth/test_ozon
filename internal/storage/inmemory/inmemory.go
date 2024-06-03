package inmemory

import (
	"sync"
	"sync/atomic"
)

type comments struct {
	storage sync.Map
	counter atomic.Uint64
}

type posts struct {
	storage sync.Map
	counter atomic.Uint64
}

type users struct {
	storage sync.Map
}

type Storage struct {
	comments comments
	posts    posts
	users    users
}

func New() *Storage {
	const op = "internal.storage.inmemory.New"
	return &Storage{}
}
