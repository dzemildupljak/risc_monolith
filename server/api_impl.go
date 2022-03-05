package main

import (
	"database/sql"

	psql "github.com/dzemildupljak/risc_monolith/server/adapters/db"
	auth_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/auth_controller"
	user_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/user_controller"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/user_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type Api struct {
	authInteractor auth_usecase.AuthInteractor
	authController auth_rest.AuthController
	userInteractor user_usecase.UserInteractor
	userController user_rest.UserController
}

func newApi(
	ac auth_rest.AuthController,
	uc user_rest.UserController) *Api {
	return &Api{
		authController: ac,
		userController: uc,
	}
}
func ApiImplementation(db sql.DB, l usecase.Logger) Api {

	pgdb := psql.New(&db)

	authRepo := psql.NewAuthRepository(*pgdb)
	userRepo := psql.NewUserRepository(*pgdb)
	authInteractor := auth_usecase.NewAuthInteractor(authRepo, l)
	userInteractor := user_usecase.NewUserInteractor(userRepo, l)
	authValidator := utils.NewAuthValidator(l)
	authController := auth_rest.NewAuthController(*authInteractor, *authValidator, l)
	userController := user_rest.NewUserController(*userInteractor, l)

	return *newApi(*authController, *userController)

}
