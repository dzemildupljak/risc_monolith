package auth_rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = "User already exists with the given email"
var ErrUserNotFound = "No user account exists with given email. Please sign in first"
var UserCreationFailed = "Unable to create user.Please try again later"

// A AuthController belong to the interface layer.
type AuthController struct {
	authInteractor auth_usecase.AuthInteractor
	logger         usecase.Logger
}

func NewAuthController(ai auth_usecase.AuthInteractor, logger usecase.Logger) *AuthController {
	return &AuthController{
		logger:         logger,
		authInteractor: ai,
	}
}
func (ac *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(auth_usecase.UserKey{}).(domain.User)

	hashedPass, err := ac.hashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: UserCreationFailed + "1"})
		return
	}

	user.Password = hashedPass
	user.Tokenhash = utils.GenerateRandomString(15)

	err = ac.authInteractor.RegisterUser(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: UserCreationFailed + "2"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&utils.GenericResponse{Status: true, Message: "user created successfully"})
}

func (ac *AuthController) hashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ac.logger.LogError("unable to hash password", "error", err)
		return "", err
	}

	return string(hashedPass), nil
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqUser := r.Context().Value(auth_usecase.UserKey{}).(domain.User)

	user, err := ac.authInteractor.UserByEmail(context.Background(), reqUser.Email)
	if err != nil {
		ac.logger.LogError("error fetching the user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to retrieve user.Please try again later"})
		return
	}
	if !user.Isverified {
		ac.logger.LogError("unverified user")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Please verify your mail address before login"})

		return
	}
	if valid := ac.authInteractor.Authenticate(&reqUser, &user); !valid {
		ac.logger.LogAccess("Authetication of user failed")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Incorrect credentials"})

		return
	}

	accessToken, err := ac.authInteractor.GenerateAccessToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to login the user. Please try again later 1"})

		return
	}

	refreshToken, err := ac.authInteractor.GenerateRefreshToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate refresh token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to login the user. Please try again later 2"})

		return
	}
	ac.logger.LogAccess("successfully generated token")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully logged in",
		Data:    &utils.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken, Username: user.Username},
	})
}

func (ac *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(auth_usecase.UserKey{}).(domain.User)
	accessToken, err := ac.authInteractor.GenerateAccessToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to generate access token.Please try again later"})

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully generated new access token",
		Data:    &utils.TokenResponse{AccessToken: accessToken},
	})
}

// Index return response which contain a listing of the resource of users.
func (uc *AuthController) Index(w http.ResponseWriter, r *http.Request) {
	uc.logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	users, err := uc.authInteractor.ShowAllUsers(r.Context())

	if err != nil {
		uc.logger.LogError("UserController-Index: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
