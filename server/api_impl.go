package main

import (
	"database/sql"

	psql "github.com/dzemildupljak/risc_monolith/server/adapers/db"
	auth_rest "github.com/dzemildupljak/risc_monolith/server/adapers/rest/auth-controller"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
)

type Api struct {
	authInteractor auth_usecase.AuthInteractor
	authController auth_rest.AuthController
}

func newApi(ai auth_usecase.AuthInteractor, ac auth_rest.AuthController) *Api {
	return &Api{
		authInteractor: ai,
		authController: ac,
	}
}
func ApiImplementation(db sql.DB, l usecase.Logger) Api {

	pgdb := psql.New(&db)

	authRepo := psql.NewAuthRepository(*pgdb)
	authInteractor := auth_usecase.NewAuthInteractor(authRepo, l)
	authController := auth_rest.NewAuthController(*authInteractor, l)

	return *newApi(*authInteractor, *authController)

}