package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/ViciousKit/course-chat-auth/models"
)

type Storage struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, name string, email string, password []byte, role int) (int64, error) {
	row := s.db.QueryRow(ctx, "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id", name, email, password, role)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, id int64) (*models.User, error) {
	row := s.db.QueryRow(ctx, "SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1", id)

	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return &models.User{}, fmt.Errorf("user with id %d not found", id)
		}

		return &models.User{}, err
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id int64, name string, email string, role int) error {
	_, err := s.db.Exec(ctx, "UPDATE users SET name = $1, email = $2, role = $3, updated_at = NOW() WHERE id = $4", name, email, role, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
