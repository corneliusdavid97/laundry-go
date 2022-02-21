package store

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/corneliusdavid97/laundry-go/src/user/service"
)

const queryGetUserByUsername = `
	select
		id,
		username,
		password,
		name,
		role
	from
		user_data
	where
		username = $1 and active = $2
	limit
		1
`

type Store struct {
	getDB func(dbName, replication string) (*sqlx.DB, error)
}

func (s *Store) GetUserByUsername(ctx context.Context, username string, active bool) (service.User, error) {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return service.User{}, err
	}
	row := db.QueryRowContext(ctx, queryGetUserByUsername, username, active)
	var user service.User
	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.Name, &user.RoleID)
	if err != nil {
		return service.User{}, err
	}
	return user, nil
}

func NewStore(getDB func(dbName, replication string) (*sqlx.DB, error)) *Store {
	return &Store{
		getDB: getDB,
	}
}
