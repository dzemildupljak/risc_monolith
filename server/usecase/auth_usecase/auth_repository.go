package auth_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type AuthRepository interface {
	DeleteUserById(ctx context.Context, id int64) error
	GetListusers(ctx context.Context) ([]domain.ShowUserParams, error)
	GetCompleteListusers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, arg domain.UpdateUserParams) (domain.User, error)
	CreateRegisterUser(ctx context.Context, arg domain.CreateRegisterUserParams) error
	GetLogedUserByEmai(ctx context.Context, username string) (domain.ShowLoginUser, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserById(ctx context.Context, id int64) (domain.User, error)
	GetBasicUserById(ctx context.Context, id int64) (domain.ShowUserParams, error)
	VerifyUserMail(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, arg domain.ChangePasswordParams) error
	GenerateResetPasswordCode(ctx context.Context, arg domain.GenerateResetPasswordCodeParams) error
	CreateOauthUser(ctx context.Context, arg domain.CreateOauthUserParams) (domain.User, error)
}

