package cache

import (
	"context"
	"social/internal/store"

	redis "github.com/go-redis/redis/v8"
)

type Storage struct {
	User interface {
		Get(ctx context.Context, userId int64) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewCacheStorage(rdb *redis.Client) Storage {
	return Storage{
		User: &UserCacheStore{rdb: rdb},
	}
}
