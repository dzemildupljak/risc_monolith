package auth_rest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/mail_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = "User already exists with the given email"
var ErrUserNotFound = "No user account exists with given email. Please sign in first"
var UserCreationFailed = "Unable to create user.Please try again later"

type ResPassword struct {
	Code                string
	Old_password        string
	New_password        string
	New_password_second string
}

// A AuthController belong to the interface layer.
type AuthController struct {
	authInteractor auth_usecase.AuthInteractor
	mailInteractor mail_usecase.MailInteractor
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

	user.MailVerfyCode, err = ac.authInteractor.RegisterUser(context.Background(), domain.CreateUserParams{Name: user.Name, Username: user.Username, Email: user.Email, Password: user.Password, Tokenhash: user.Tokenhash})
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: UserCreationFailed + "2"})
		return
	}

	verifyMail := mail_usecase.Mail{
		Reciever:  user.Email,
		MailTitle: "Verify email",
		Type:      1,
	}

	ac.mailInteractor.SendEmail(verifyMail, user.MailVerfyCode, user.Name)

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

func (ac *AuthController) VerifyMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ac.logger.LogAccess("verifying the confimation code")
	verificationData := r.Context().Value(VerificationDataKey{}).(VerificationData)

	err := ac.authInteractor.AuthRepository.VerifyUserMail(r.Context(), verificationData.Email)
	if err != nil {
		ac.logger.LogError("Failed to verify user mail try again later", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Failed to verify user mail try again later"})
	}

	ac.logger.LogAccess("successfully verified mail")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "successfully verified mail",
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
func (ac *AuthController) PasswordResetCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(auth_usecase.UserIDKey{}).(string)
	usrId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ac.logger.LogError("code validation failed", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Generating failed. Password reset code malformed"})
		return
	}

	user, err := ac.authInteractor.UserById(context.Background(), usrId)
	if err != nil {
		ac.logger.LogError("unable to get user to generate secret code for password reset", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to send password reset code. Please try again later1"})
		return
	}

	resPassData := &ResPassword{}

	err = json.NewDecoder(r.Body).Decode(resPassData)
	if err != nil {
		ac.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	hashOldPass, err := ac.hashPassword(resPassData.Old_password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to reset password. Please try again later2"})
		return
	}

	fmt.Println("hashOldPass", hashOldPass)
	fmt.Println("user.Password", user.Password)
	fmt.Println("resPassData.New_password", resPassData.New_password)
	fmt.Println("resPassData.New_password_second", resPassData.New_password_second)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(resPassData.Old_password))
	if err != nil {
		ac.authInteractor.Logger.LogError("old password are not same")
	}

	if resPassData.New_password != resPassData.New_password_second {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to reset password. Please try again later3"})
		return
	}

	hashNewPass, err := ac.hashPassword(resPassData.New_password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to reset password. Please try again later4"})
		return
	}
	err = ac.authInteractor.AuthRepository.ChangePassword(r.Context(), domain.ChangePasswordParams{Email: user.Email, Password: hashNewPass})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to reset password. Please try again later"})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully reseted passowrd",
		Data:    nil,
	})

}

func (ac *AuthController) GeneratePassResetCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(auth_usecase.UserIDKey{}).(string)
	usrId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ac.logger.LogError("code validation failed", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Generating failed. Password reset code malformed"})
		return
	}

	user, err := ac.authInteractor.UserById(context.Background(), usrId)
	if err != nil {
		ac.logger.LogError("unable to get user to generate secret code for password reset", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Unable to send password reset code. Please try again later"})
		return
	}

	verfyPassword := mail_usecase.Mail{
		Reciever:  user.Email,
		MailTitle: "Password reset code",
		Type:      2,
	}

	//TODO geenrate password reset code and query to write password code into table
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	user.PasswordVerfyCode = fmt.Sprint(rand.Intn(max-min+1) + min)
	user.PasswordVerfyExpire = time.Now().Add(1 * time.Hour)

	ac.mailInteractor.SendEmail(verfyPassword, user.PasswordVerfyCode, user.Name)

	ac.logger.LogAccess("successfully mailed password reset code")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Please check your mail for password reset code"})
}

// Index return response which contain a listing of the resource of users.
func (uc *AuthController) Index(w http.ResponseWriter, r *http.Request) {
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
