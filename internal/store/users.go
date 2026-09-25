package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO USERS (username, email, password) 
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) // what about retrying if timeout does occurs
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at 
		FROM USERS 
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		// create user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create user invitation
		if err := s.createUserInvite(ctx, tx, user.Id, token, exp); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// 1. check user exists
	// 2. update/activate the user
	// 3. delete the token after activation
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		var user User
		if err := s.getUserFromInvitation(ctx, tx, token, &user); err != nil {
			return err
		}

		user.IsActive = true

		if err := s.Update(ctx, tx, &user); err != nil {
			return err
		}

		if err := s.deleteUserToken(ctx, tx, user.Id); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE USERS 
		SET username = $2, email = $3, is_active = $4
		WHERE id = $1 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Id, user.Username, user.Email, user.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userId int64) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		// Delete the user
		if err := s.deleteUser(ctx, tx, userId); err != nil {
			return err
		}

		if err := s.deleteUserToken(ctx, tx, userId); err != nil {
			return err
		}
		// Delete the invite token
		return nil
	})
}

func (s *UserStore) deleteUser(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `
		DELETE FROM USERS WHERE id = $1 
	`

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserToken(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `
		DELETE FROM USER_INVITATIONS WHERE user_id = $1
	`

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string, user *User) error {
	query := `
			SELECT u.id, u.username, u.email, u.created_at, u.is_active 
			FROM USERS u 
			JOIN USER_INVITATIONS ui ON u.id = ui.user_id
			WHERE ui.token = $1 AND ui.expiry > $2
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, token, time.Now()).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) createUserInvite(ctx context.Context, tx *sql.Tx, userId int64, token string, exp time.Duration) error {
	query := `
		INSERT INTO USER_INVITATIONS 
		(token, user_id, expiry) 
		VALUES 
		($1, $2, $3) 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}
