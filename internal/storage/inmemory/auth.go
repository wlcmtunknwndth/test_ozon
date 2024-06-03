package inmemory

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/test_ozon/graph/model"
)

//var usersCount uint64 = 0

func (s *Storage) CreateUser(ctx context.Context, usr *model.NewUser) error {
	const op = "internal.storage.inmemory.auth.CreateUser"
	s.users.storage.Store(usr.Username, model.User{
		Username: usr.Username,
		Password: usr.Password,
		IsAdmin:  false,
	})
	return nil
}

func (s *Storage) GetPassword(ctx context.Context, username string) (string, error) {
	const op = "internal.storage.inmemory.auth.GetPassword"
	data, ok := s.users.storage.Load(username)
	if !ok {
		return "", fmt.Errorf("%s: couldn't find value", op)
	}
	usr, ok := data.(model.User)
	if !ok {
		return "", fmt.Errorf("%s: couldn't transform value into user", op)
	}
	return usr.Password, nil
}

func (s *Storage) IsAdmin(ctx context.Context, username string) (bool, error) {
	const op = "internal.storage.inmemory.auth.IsAdmin"
	data, ok := s.users.storage.Load(username)
	if !ok {
		return false, fmt.Errorf("%s: couldn't find value", op)
	}
	usr, ok := data.(model.User)
	if !ok {
		return false, fmt.Errorf("%s: couldn't transform value into user", op)
	}
	return usr.IsAdmin, nil
}

func (s *Storage) DeleteUser(ctx context.Context, username string) error {
	const op = "internal.storage.inmemory.auth.DeleteUser"
	s.users.storage.Delete(username)
	return nil
}
