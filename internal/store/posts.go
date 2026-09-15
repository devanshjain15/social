package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserId    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int64     `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int64 `json:"comment_count"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.UserId).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt) // what is this scan doing and when we are creating a new post using the Post struct where does the created and updated timestamp values are obtained does scan does it?
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	var post Post
	query := `
		SELECT id, title, content, user_id, created_at, tags, updated_at, version FROM POSTS 
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, postId).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, pq.Array(&post.Tags), &post.UpdatedAt, &post.Version)
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

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM POSTS WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE POSTS 
		SET title = $2, content = $3, version = version + 1
		WHERE id = $1 AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Id, post.Title, post.Content, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.content,
			p.tags,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count
		FROM
			posts p
			LEFT JOIN comments c ON p.id = c.post_id
			LEFT JOIN users u ON p.user_id = u.id
		WHERE
			(
				p.user_id = $1
				OR EXISTS (
					SELECT
						1
					FROM
						followers f
					WHERE
						f.user_id = p.user_id
						AND f.follower_id = $1
				)
			)
			AND (
				$5 = ''
				OR p.title ILIKE '%' || $5 || '%'
				OR p.content ILIKE '%' || $5 || '%'
			)
			AND (
				$4::varchar[] IS NULL
				OR $4::varchar[] = '{}'
				OR p.tags @> $4::varchar[]
			)
		GROUP BY
			p.id,
			u.id
		ORDER BY
			p.created_at ` + fq.Sort + `
		LIMIT
			$2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, pq.Array(fq.Tags), fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	feed := make([]PostWithMetaData, 0)

	for rows.Next() {
		p := PostWithMetaData{}
		if err := rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			pq.Array(&p.Tags),
			&p.UserId,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.User.Username,
			&p.CommentCount,
		); err != nil {
			return nil, err
		}

		p.User.Id = p.UserId

		feed = append(feed, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feed, nil
}
