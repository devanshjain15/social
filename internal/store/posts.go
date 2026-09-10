package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	UserId    int64    `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Comments []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB // why passing pointer to the database? why struct this way?
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO POSTS (title, content, tags, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.UserId).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt) // what is this scan doing and when we are creating a new post using the Post struct where does the created and updated timestamp values are obtained does scan does it?
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	var post Post
	query := `
		SELECT id, title, content, user_id, created_at, tags, updated_at FROM POSTS 
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, query, postId).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, pq.Array(&post.Tags), &post.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}
