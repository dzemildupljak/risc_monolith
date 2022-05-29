package user_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
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

type UserInteractor struct {
	userRepository UserRepository
	Logger         usecase.Logger
}

func NewUserInteractor(r UserRepository, l usecase.Logger) *UserInteractor {
	return &UserInteractor{
		userRepository: r,
		Logger:         l,
	}
}

func (ui *UserInteractor) ListUsersInteract(
	ctx context.Context) ([]domain.ShowUserParams, error) {

	users, err := ui.userRepository.UserList(ctx)
	return users, err
}

func (ui *UserInteractor) ListCompleteUsersInteract(
	ctx context.Context) ([]domain.User, error) {

	users, err := ui.userRepository.ListCompleteUsers(ctx)
	return users, err
}

func (ui *UserInteractor) UserByIdInteract(
	ctx context.Context, userId int64) (domain.ShowUserParams, error) {

	user, err := ui.userRepository.UserById(ctx, userId)
	return user, err
}

func (ui *UserInteractor) UserByEmailInteract(
	ctx context.Context, userEmail string) (domain.ShowUserParams, error) {

	user, err := ui.userRepository.UserByEmail(ctx, userEmail)
	return user, err
}

func (ui *UserInteractor) UserUpdate(
	ctx context.Context, uId int64, usr domain.UpdateUserParams) (domain.ShowUserParams, error) {

	user, err := ui.userRepository.UpdateUser(ctx, uId, usr)

	return user, err
}
