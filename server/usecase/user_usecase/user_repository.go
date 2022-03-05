package user_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type UserRepository interface {
	UserList(ctx context.Context) ([]domain.ShowUserParams, error)
	ListCompleteUsers(ctx context.Context) ([]domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.ShowUserParams, error)
	CompleteUserById(ctx context.Context, id int64) (domain.User, error)
	UserById(ctx context.Context, id int64) (domain.ShowUserParams, error)
	UserDeleteById(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, usrId int64, arg domain.UpdateUserParams) (domain.ShowUserParams, error)
}

// TODO add patch and put (update) for user CRUD
