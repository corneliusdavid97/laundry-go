package user

import (
	"context"
	"errors"
)

type User struct {
	UserID   int64
	Username string
	Name     string
	Role     Role
}

var ErrAuthFailed = errors.New("Username atau password salah")

type Service interface {
	AuthUser(ctx context.Context, username, password string) (User, error)
}

var defaultService Service

func Init(s Service) {
	defaultService = s
}

func GetService() Service {
	return defaultService
}
