package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ViciousKit/course-chat-auth/models"
)

type Storage struct {
	db *sql.DB
}

func New(user string, password string, dbname string, host string, port int) *Storage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Cant connect pg" + err.Error())
		panic(err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Cant ping pg" + err.Error())
		panic(err)
	}
	fmt.Println("Connected!")

	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, name string, email string, password []byte, role int) error {
	method := "CreateUser"

	statement, err := s.db.Prepare("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	_, err = statement.ExecContext(ctx, name, email, password, role)
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, id int64) (*models.User, error) {
	method := "GetUser"

	statement, err := s.db.Prepare("SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}

	var user models.User
	err = statement.QueryRowContext(ctx, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.User{}, fmt.Errorf("%s: %s", method, "user not found")
		}

		return &models.User{}, fmt.Errorf("%s: %w", method, err)
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id int64, name string, email string, role int) error {
	method := "UpdateUser"

	statement, err := s.db.Prepare("UPDATE users SET name = $1, email = $2, role = $3, updated_at = NOW() WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	_, err = statement.ExecContext(ctx, name, email, role, id)
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int64) error {
	method := "DeleteUser"

	statement, err := s.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	_, err = statement.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}

	return nil
}
