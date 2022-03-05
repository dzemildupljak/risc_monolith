package user_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

type UserInteractor struct {
	UserRepository UserRepository
	Logger         usecase.Logger
}

func NewUserInteractor(r UserRepository, l usecase.Logger) *UserInteractor {
	return &UserInteractor{
		UserRepository: r,
		Logger:         l,
	}
}

func (ui *UserInteractor) ListUsersInteract(
	ctx context.Context) ([]domain.ShowUserParams, error) {

	users, err := ui.UserRepository.UserList(ctx)
	return users, err
}

func (ui *UserInteractor) ListCompleteUsersInteract(
	ctx context.Context) ([]domain.User, error) {

	users, err := ui.UserRepository.ListCompleteUsers(ctx)
	return users, err
}

func (ui *UserInteractor) UserByIdInteract(
	ctx context.Context, userId int64) (domain.ShowUserParams, error) {

	user, err := ui.UserRepository.UserById(ctx, userId)
	return user, err
}

func (ui *UserInteractor) UserByEmailInteract(
	ctx context.Context, userEmail string) (domain.ShowUserParams, error) {

	user, err := ui.UserRepository.UserByEmail(ctx, userEmail)
	return user, err
}

func (ui *UserInteractor) UserUpdate(
	ctx context.Context, uId int64, usr domain.UpdateUserParams) (domain.ShowUserParams, error) {

	user, err := ui.UserRepository.UpdateUser(ctx, uId, usr)

	return user, err
}
