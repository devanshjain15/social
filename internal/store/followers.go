package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FollowerStore struct {
	db *sql.DB
}

type Follower struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (s *FollowerStore) Follow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		INSERT INTO FOLLOWERS (user_id, follower_id) VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return nil

}

// idempotency
func (s *FollowerStore) Unfollow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		DELETE FROM FOLLOWERS WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		return err
	}
	return nil
}
