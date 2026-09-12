package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	PostId    int64  `json:"post_id"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO COMMENTS (content, post_id, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.Content, comment.PostId, comment.UserId).Scan(&comment.Id, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, u.username, u.id, u.created_at FROM COMMENTS c 
		JOIN USERS u ON c.user_id = u.id
		WHERE c.post_id = $1 
		ORDER BY c.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, 0)

	for rows.Next() {
		var comment Comment
		var user User

		if err := rows.Scan(
			&comment.Id,
			&comment.Content,
			&comment.PostId,
			&comment.UserId,
			&comment.CreatedAt,
			&user.Username,
			&user.Id,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		comment.User = user
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
