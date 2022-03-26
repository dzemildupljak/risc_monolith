package auth_rest

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase/mail_usecase"
)

type AuthUsecase interface {
	Authenticate(reqUser *domain.User, user *domain.User) bool
	GenerateAccessToken(user *domain.User) (string, error)
	ValidateAccessToken(tokenString string) (string, string, error)
	GenerateCustomKey(userID string, tokenHash string) string
	GenerateRefreshToken(user *domain.User) (string, error)
	ValidateRefreshToken(tokenString string) (string, string, error)
	RegisterUser(
		ctx context.Context, u domain.CreateUserParams) (string, error)
	RegisterOauthUser(
		ctx context.Context, u domain.CreateOauthUserParams) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, error)
	UserById(ctx context.Context, usrID int64) (domain.User, error)
	BasicUserById(
		ctx context.Context, usrID int64) (domain.ShowUserParams, error)
	ShowAllUsers(ctx context.Context) ([]domain.ShowUserParams, error)
	ShowCompleteUsers(ctx context.Context) ([]domain.User, error)
	GenerateResetPasswCode(
		ctx context.Context, email string) (mail_usecase.Mail, string, error)
	UserMailVerify(ctx context.Context, email string) error
	UpdatePassword(
		ctx context.Context, usr domain.ChangePasswordParams) error
}
