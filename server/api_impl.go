package main

import (
	"database/sql"

	psql "github.com/dzemildupljak/risc_monolith/server/adapters/db"
	auth_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/auth_controller"
	donor_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/donor_controller"
	user_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/user_controller"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/donor_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/user_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type Api struct {
	authController  auth_rest.AuthController
	userController  user_rest.UserController
	donorController donor_rest.DonorController
}

func newApi(
	ac auth_rest.AuthController,
	uc user_rest.UserController,
	dc donor_rest.DonorController) *Api {
	return &Api{
		authController:  ac,
		userController:  uc,
		donorController: dc,
	}
}
func ApiImplementation(db sql.DB, l usecase.Logger) Api {

	pgdb := psql.New(&db)

	authRepo := psql.NewAuthRepository(*pgdb)
	userRepo := psql.NewUserRepository(*pgdb)
	donorRepo := psql.NewDonorRepository(*pgdb)

	mapper := utils.NewMapper()

	authInteractor := auth_usecase.NewAuthInteractor(authRepo, l)
	userInteractor := user_usecase.NewUserInteractor(userRepo, l)
	donorInteractor := donor_usecase.NewDonorInteractor(donorRepo, l, *mapper)

	authValidator := utils.NewAuthValidator(l)

	authController := auth_rest.NewAuthController(*authInteractor, *authValidator, l)
	userController := user_rest.NewUserController(*userInteractor, l)
	donorController := donor_rest.NewDonorController(*donorInteractor, l)

	return *newApi(*authController, *userController, *donorController)

}
