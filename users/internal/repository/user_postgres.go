package repository

import (
	"awesomeProject1/internal/entity"
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) (int64, error) {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash)

	if err := row.Scan(&user.ID); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := entity.User{}
	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
