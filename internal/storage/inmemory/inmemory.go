package inmemory

import (
	"github.com/arriqaaq/flashdb"
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

func New(config *flashdb.Config) *Storage {
	const op = "internal.storage.inmemory.New"
	return &Storage{}
}
