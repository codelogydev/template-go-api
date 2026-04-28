package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/codelogydev/template-go-api/internal/model"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
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
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}

	return users, nil
}
