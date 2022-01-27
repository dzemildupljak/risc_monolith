package auth_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error)
	DeleteUserById(ctx context.Context, id int64) error
	GetListusers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, arg domain.UpdateUserParams) (domain.User, error)
	CreateRegisterUser(ctx context.Context, arg domain.CreateRegisterUserParams) error
	GetUserByUsername(ctx context.Context, username string) (domain.ShowLoginUser, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserById(ctx context.Context, id int64) (domain.User, error)
}
