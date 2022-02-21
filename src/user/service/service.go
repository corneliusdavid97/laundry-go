package service

import (
	"context"
	"database/sql"

	"github.com/corneliusdavid97/laundry-go/src/user"
)

type User struct {
	UserID   int64
	Username string
	Name     string
	RoleID   int
	Password string
}

type Service struct {
	store Store
}

type Store interface {
	GetUserByUsername(ctx context.Context, username string, active bool) (User, error)
}

func (s *Service) AuthUser(ctx context.Context, username, password string) (user.User, error) {
	userTmp, err := s.store.GetUserByUsername(ctx, username, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrAuthFailed
		}
		return user.User{}, err
	}

	if userTmp.Password != password {
		return user.User{}, user.ErrAuthFailed
	}

	return parseUser(userTmp), nil
}

func parseUser(u User) user.User {
	return user.User{
		UserID:   u.UserID,
		Name:     u.Name,
		Username: u.Username,
		Role: user.Role{
			RoleID:   user.RoleID(u.RoleID),
			RoleName: user.GetRoleNameByRoleID(user.RoleID(u.RoleID)),
		},
	}
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}
