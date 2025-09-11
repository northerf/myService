package repository

import (
	"awesomeProject1/users/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type Repository struct {
	UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepo(db),
	}
}
