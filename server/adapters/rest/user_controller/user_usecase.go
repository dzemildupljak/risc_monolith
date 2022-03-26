package user_rest

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type UserUsecaseRepository interface {
	ListUsersInteract(ctx context.Context) ([]domain.ShowUserParams, error)
	ListCompleteUsersInteract(ctx context.Context) ([]domain.User, error)
	UserByIdInteract(ctx context.Context, userId int64) (domain.ShowUserParams, error)
	UserByEmailInteract(
		ctx context.Context, userEmail string) (domain.ShowUserParams, error)
	UserUpdate(
		ctx context.Context, uId int64,
		usr domain.UpdateUserParams) (domain.ShowUserParams, error)
}
