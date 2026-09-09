package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Record not found!")
)

type Storage struct {
	Posts interface { // why interface though and if it is being used at multiple places why not create a type out of it?
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
	}

	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
