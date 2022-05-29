package main

import (
	"database/sql"

	"github.com/dzemildupljak/risc_monolith/server/adapters/db/psql"
	auth_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/auth_controller"
	donation_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/donation_controller"
	donor_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/donor_controller"
	user_rest "github.com/dzemildupljak/risc_monolith/server/adapters/rest/user_controller"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	donationevent_usecase "github.com/dzemildupljak/risc_monolith/server/usecase/donation_event"
	"github.com/dzemildupljak/risc_monolith/server/usecase/donor_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/user_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type Api struct {
	authController     auth_rest.AuthController
	userController     user_rest.UserController
	donorController    donor_rest.DonorController
	donationController donation_rest.DonationEventController
}

func newApi(
	ac auth_rest.AuthController,
	uc user_rest.UserController,
	dc donor_rest.DonorController,
	dec donation_rest.DonationEventController,
) *Api {
	return &Api{
		authController:     ac,
		userController:     uc,
		donorController:    dc,
		donationController: dec,
	}
}
func ApiImplementation(db sql.DB, l usecase.Logger) Api {

	pgdb := psql.New(&db)

	authRepo := psql.NewAuthRepository(*pgdb)
	userRepo := psql.NewUserRepository(*pgdb)
	donorRepo := psql.NewDonorRepository(*pgdb)
	donationRepo := psql.NewDonationEventRepository(*pgdb)

	mapper := utils.NewMapper()

	authInteractor := auth_usecase.NewAuthInteractor(authRepo, l)
	userInteractor := user_usecase.NewUserInteractor(userRepo, l)
	donorInteractor := donor_usecase.NewDonorInteractor(donorRepo, l, *mapper)
	donationInteractor := donationevent_usecase.NewDonationEventInteractor(donationRepo, l)

	authValidator := utils.NewAuthValidator(l)

	authController := auth_rest.NewAuthController(authInteractor, *authValidator, l)
	userController := user_rest.NewUserController(userInteractor, l)
	donorController := donor_rest.NewDonorController(donorInteractor, l)
	donationController := donation_rest.NewDonationEventController(donationInteractor, l)

	return *newApi(*authController, *userController, *donorController, *donationController)

}
