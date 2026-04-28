package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/codelogydev/template-go-api/internal/model"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, name, email string) (*model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]model.User, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email",
		name, email).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}


