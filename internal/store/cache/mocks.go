package cache

import (
	"context"
	"social/internal/store"
)

func NewMockCacheStorage() Storage {
	return Storage{
		User: &MockUserCache{},
	}
}

type MockUserCache struct{}

func (cache *MockUserCache) Get(ctx context.Context, userId int64) (*store.User, error) {
	return nil, nil
}
func (cache *MockUserCache) Set(ctx context.Context, user *store.User) error {
	return nil
}
