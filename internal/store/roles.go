package store

import (
	"context"
	"database/sql"
)

type RoleStore struct {
	db *sql.DB
}

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

func (s *RoleStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `SELECT id, name, level, description FROM ROLES WHERE name = $1`
	role := Role{}
	err := s.db.QueryRowContext(ctx, query, roleName).Scan(&role.Id, &role.Name, &role.Level, &role.Description)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
