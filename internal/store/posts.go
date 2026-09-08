package store

import (
	"context"
	"database/sql"

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
