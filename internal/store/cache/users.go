package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"social/internal/store"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const ttl = time.Minute

type UserCacheStore struct {
	rdb *redis.Client
}

func (cache *UserCacheStore) Get(ctx context.Context, userId int64) (*store.User, error) {
	key := fmt.Sprintf("user:%v", userId)
	data, err := cache.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (cache *UserCacheStore) Set(ctx context.Context, user *store.User) error {
	key := fmt.Sprintf("user:%v", user.Id)
	userString, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.rdb.SetEX(ctx, key, userString, ttl).Err()
}
