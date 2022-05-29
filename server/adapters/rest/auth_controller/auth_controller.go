package auth_rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/mail_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type ResPassword struct {
	Code                string
	Old_password        string
	New_password        string
	New_password_second string
}

type SetPasswordValues struct {
	New_password        string `json:"new_password"`
	New_password_second string `json:"new_password_second"`
}

// A AuthController belong to the interface layer.
type AuthController struct {
	ai     auth_usecase.AuthUsecase
	av     utils.AuthValidator
	mi     mail_usecase.MailUsecase
	logger usecase.Logger
}

func NewAuthController(
	ai auth_usecase.AuthUsecase,
	av utils.AuthValidator,
	logger usecase.Logger) *AuthController {
	return &AuthController{
		logger: logger,
		ai:     ai,
		av:     av,
	}
}

func (ac *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := &domain.CreateRegisterUserParams{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		ac.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create user.Please try again later",
			})
		return
	}

	hashedPass, err := ac.hashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create user.Please try again later",
			})
		return
	}

	user.Password = hashedPass
	user.Tokenhash = utils.GenerateRandomString(15)

	user.MailVerfyCode, err = ac.ai.RegisterUser(
		context.Background(),
		domain.CreateUserParams{
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			Tokenhash: user.Tokenhash,
		})

	if err != nil {
		ac.logger.LogError("err", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create user.Please try again later",
			})
		return
	}

	verifyMail := mail_usecase.Mail{
		Reciever:  user.Email,
		MailTitle: "Verify email",
		Type:      1,
	}

	ac.mi.SendEmail(verifyMail, user.MailVerfyCode, user.Name)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "User created successfully",
		})
}

func (ac *AuthController) hashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ac.logger.LogError("unable to hash password", "error", err)
		return "", err
	}

	return string(hashedPass), nil
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Decode/parse struct from request
	logedUser := &domain.ShowLoginUser{}
	err := json.NewDecoder(r.Body).Decode(logedUser)
	if err != nil {

		ac.logger.LogError("login1 = deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Message: "Invalid credentials"})
		return
	}

	// Validate login values ex. email/password != ""
	errRes, err := ac.av.ValidateLoginValues(*logedUser)
	if err != nil {

		ac.logger.LogError("login2 = deserialization of user json failed", "error", err)

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Data:    errRes,
			Message: "Invalid credentials"})
		return
	}

	// Get user by email from interactor-repository
	user, err := ac.ai.UserByEmail(
		context.Background(), logedUser.Email)

	if err != nil {
		ac.logger.LogError("error fetching the user", "error", err)

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "No user account exists with given email. Please sign up first",
			})
		return
	}
	// Check if given user email is verified
	if !user.Isverified {
		ac.logger.LogError("unverified user")

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Please verify your mail address before login",
			})

		return
	}

	// Check if given password is same is password in db(hashed)
	if valid := ac.ai.Authenticate(
		&domain.User{
			Email:    logedUser.Email,
			Password: logedUser.Password,
		},
		&user); !valid {
		ac.logger.LogAccess("Authetication of user failed")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Incorrect credentials",
			})

		return
	}

	// Generate access jwt token with payload and signature
	accessToken, err := ac.ai.GenerateAccessToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to login the user. Please try again later",
			})

		return
	}

	// Generate refresh jwt token with payload and signature
	refreshToken, err := ac.ai.GenerateRefreshToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate refresh token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to login the user. Please try again later",
			})

		return
	}
	ac.logger.LogAccess("successfully generated token")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully logged in",
		Data: &utils.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Email:        user.Email,
		},
	})
}

func (ac *AuthController) VerifyMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ac.logger.LogAccess("verifying the confimation code")
	verificationData := r.Context().Value(VerificationDataKey{}).(VerificationData)

	err := ac.ai.UserMailVerify(r.Context(), verificationData.Email)
	if err != nil {
		ac.logger.LogError("Failed to verify user mail try again later", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Failed to verify user mail try again later",
			})
	}

	ac.logger.LogAccess("successfully verified mail")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully verified mail",
	})
}

func (ac *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(auth_usecase.UserKey{}).(domain.User)
	accessToken, err := ac.ai.GenerateAccessToken(&user)
	if err != nil {
		ac.logger.LogError("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to generate access token.Please try again later",
			})

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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	user, err := ac.ai.UserById(context.Background(), usrId)
	if err != nil {
		ac.logger.LogError(
			"unable to get user to generate secret code for password reset", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to send password reset code. Please try again later",
			})
		return
	}

	resPassData := &ResPassword{}

	err = json.NewDecoder(r.Body).Decode(resPassData)
	if err != nil {
		ac.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	validateExTime := utils.ValidateExpirationTime(user.PasswordVerfyExpire.Time)

	if !validateExTime {
		ac.logger.LogError("verification code failed2")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(resPassData.Old_password))

	if err != nil {
		ac.logger.LogError("old password invalid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	if user.PasswordVerfyCode != resPassData.Code {
		ac.logger.LogError(
			"requested code and user code are not same - reset password")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	if resPassData.New_password != resPassData.New_password_second {
		ac.logger.LogError(
			"new_password and new_password_second are different - reset password")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	hashNewPass, err := ac.hashPassword(resPassData.New_password)
	if err != nil {
		ac.logger.LogError(
			"hasing new password failed - reset password")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	err = ac.ai.UpdatePassword(
		r.Context(),
		domain.ChangePasswordParams{
			Email:    user.Email,
			Password: hashNewPass,
		})
	if err != nil {
		ac.logger.LogError(
			"updateting user password faield - reset password")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully reseted password",
		Data:    nil,
	})
}

func (ac *AuthController) GeneratePassResetCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(auth_usecase.UserIDKey{}).(string)
	usrId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ac.logger.LogError("code validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Generating failed. Password reset code malformed",
			})
		return
	}

	user, err := ac.ai.UserById(context.Background(), usrId)
	if err != nil {
		ac.logger.LogError(
			"unable to get user to generate code for password reset", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to send password reset code. Please try again later",
			})
		return
	}

	verfyPassMail, verfyPassCode, err := ac.ai.GenerateResetPasswCode(
		r.Context(),
		user.Email,
	)

	if err != nil {
		ac.logger.LogError("unable to get user to reset password", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to send reset password. Please try again later",
			})
		return
	}

	ac.mi.SendEmail(verfyPassMail, verfyPassCode, user.Name)

	ac.logger.LogAccess("successfully mailed password reset code")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "Please check your mail for password reset code",
		})
}

func (ac *AuthController) SetNewPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(auth_usecase.UserIDKey{}).(string)
	usrId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ac.logger.LogError("code validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Generating failed. Password malformed",
			})
		return
	}
	settPassData := &SetPasswordValues{}

	err = json.NewDecoder(r.Body).Decode(settPassData)
	if err != nil {
		ac.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	if settPassData.New_password != settPassData.New_password_second {
		ac.logger.LogError(
			"unable to get user to set password different pass values", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to set new password. Please try again later",
			})
		return
	}

	user, err := ac.ai.UserById(context.Background(), usrId)
	if err != nil {
		ac.logger.LogError(
			"unable to get user to generate code for password reset", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to set new password. Please try again later",
			})
		return
	}
	if user.Password != "" {
		ac.logger.LogError(
			"user already have password", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to set new password. Please try again later",
			})
		return
	}
	hashNewPass, err := ac.hashPassword(settPassData.New_password)
	if err != nil {
		ac.logger.LogError(
			"hasing new password failed - reset password")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	err = ac.ai.UpdatePassword(
		r.Context(),
		domain.ChangePasswordParams{
			Email:    user.Email,
			Password: hashNewPass,
		},
	)
	if err != nil {
		ac.logger.LogError(
			"unable to set user to new password in db", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to set new password. Please try again later",
			})
		return
	}

	ac.logger.LogAccess("successfully setted new password")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "successfully set new password",
		})

}

// Index return response which contain a listing of the resource of users.
func (uc *AuthController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := uc.ai.ShowCompleteUsers(r.Context())

	if err != nil {
		uc.logger.LogError("UserController-Index: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Index return response which contain a listing of the resource of users.
func (uc *AuthController) UserIndex(w http.ResponseWriter, r *http.Request) {
	users, err := uc.ai.ShowAllUsers(r.Context())

	if err != nil {
		uc.logger.LogError("UserController-Index: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserById return response of the resource of users.
func (uc *AuthController) BasicUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["user_id"], 10, 64)
	if err != nil {
		uc.logger.LogError("BasicUserById = user Id validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get user. Please try again later",
			})
		return
	}
	usr, err := uc.ai.BasicUserById(r.Context(), userId)

	if err != nil {
		uc.logger.LogError("get basic user by id", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get user. Please try again later",
			})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// ForgotPassword - after some validation with code, email, expiration
// update hashed password for user with given email
func (ac *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqUser := &domain.ForgotPasswordValues{}
	err := json.NewDecoder(r.Body).Decode(reqUser)
	if err != nil {

		ac.logger.LogError("deserialization of user json failed1", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Message: "Invalid credentials"})
		return
	}
	errRes, err := ac.av.ValidateForgotPassValues(*reqUser)

	if err != nil {

		ac.logger.LogError("deserialization of user json failed2", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Data:    errRes,
			Message: "Invalid credentials"})
		return
	}

	usr, err := ac.ai.UserByEmail(r.Context(), reqUser.Email)

	if err != nil {
		ac.logger.LogError("no user with atempted email")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	if usr.PasswordVerfyCode != reqUser.Code {
		ac.logger.LogError("verify code and req code are different")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}
	if usr.Email != reqUser.Email {
		ac.logger.LogError("verify email and req email are different")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	if !utils.ValidateExpirationTime(usr.PasswordVerfyExpire.Time) {
		ac.logger.LogError("forgot password code is expired")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	hashNewPass, err := ac.hashPassword(reqUser.New_password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to reset password. Please try again later",
			})
		return
	}

	err = ac.ai.UpdatePassword(
		r.Context(),
		domain.ChangePasswordParams{
			Email:    reqUser.Email,
			Password: hashNewPass,
		})

	if err != nil {
		ac.logger.LogError("deserialization of user json failed3", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Message: "Invalid credentials"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&utils.GenericResponse{
		Status:  true,
		Message: "Successfully reseted password",
		Data:    nil,
	})
}

func (ac *AuthController) ForgotPasswordCode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqUser := &domain.UserEmail{}
	err := json.NewDecoder(r.Body).Decode(reqUser)
	if err != nil {

		ac.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&utils.GenericResponse{
			Status:  false,
			Message: "Invalid credentials"})
		return
	}
	err = validEmail(reqUser.Email)

	if err != nil {
		ac.logger.LogError("inalid email address", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Invalid credentials",
			})
		return
	}

	user, err := ac.ai.UserByEmail(context.Background(), reqUser.Email)
	if err != nil {
		ac.logger.LogError(
			"unable to get user to generate code for password reset", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to send code. Please try again later",
			})
		return
	}

	verfyPassMail, verfyPassCode, err := ac.ai.GenerateResetPasswCode(
		r.Context(),
		user.Email,
	)

	if err != nil {
		ac.logger.LogError("unable to get user to reset password", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to send code. Please try again later",
			})
		return
	}

	ac.mi.SendEmail(verfyPassMail, verfyPassCode, user.Name)

	ac.logger.LogAccess("successfully mailed password reset code")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "Please check your mail for code",
		})
}
