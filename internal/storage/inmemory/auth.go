package inmemory

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
)

var usersCount uint64 = 0

func (s *Storage) CreateUser(ctx context.Context, usr *model.NewUser) error {
	const op = "internal.storage.inmemory.auth.CreateUser"
	s.users.storage[usr.Username] = model.User{
		Username: usr.Username,
		Password: usr.Password,
		IsAdmin:  false,
	}
	return nil
}

func (s *Storage) GetPassword(ctx context.Context, username string) (string, error) {
	const op = "internal.storage.inmemory.auth.GetPassword"
	return s.users.storage[username].Username, nil
}

func (s *Storage) IsAdmin(ctx context.Context, username string) (bool, error) {
	const op = "internal.storage.inmemory.auth.IsAdmin"
	return s.users.storage[username].IsAdmin, nil
}

func (s *Storage) DeleteUser(ctx context.Context, username string) error {
	const op = "internal.storage.inmemory.auth.DeleteUser"
	delete(s.users.storage, username)
	return nil
}
