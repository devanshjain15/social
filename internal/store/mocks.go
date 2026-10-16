package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStorage() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (s *MockUserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	return nil, nil
}
func (s *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
func (s *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}
func (s *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}
func (s *MockUserStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}
func (s *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (s *MockUserStore) Delete(ctx context.Context, userId int64) error {
	return nil
}
