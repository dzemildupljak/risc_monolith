package auth_usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/mail_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// UserKey is used as a key for storing the User object in context at middleware
type UserKey struct{}

// UserIDKey is used as a key for storing the UserID in context at middleware
type UserIDKey struct{}

type AuthUsecase interface {
	Authenticate(reqUser *domain.User, user *domain.User) bool
	GenerateAccessToken(user *domain.User) (string, error)
	ValidateAccessToken(tokenString string) (string, string, error)
	GenerateCustomKey(userID string, tokenHash string) string
	GenerateRefreshToken(user *domain.User) (string, error)
	ValidateRefreshToken(tokenString string) (string, string, error)
	RegisterUser(
		ctx context.Context, u domain.CreateUserParams) (string, error)
	RegisterOauthUser(
		ctx context.Context, u domain.CreateOauthUserParams) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, error)
	UserById(ctx context.Context, usrID int64) (domain.User, error)
	BasicUserById(
		ctx context.Context, usrID int64) (domain.ShowUserParams, error)
	ShowAllUsers(ctx context.Context) ([]domain.ShowUserParams, error)
	ShowCompleteUsers(ctx context.Context) ([]domain.User, error)
	GenerateResetPasswCode(
		ctx context.Context, email string) (mail_usecase.Mail, string, error)
	UserMailVerify(ctx context.Context, email string) error
	UpdatePassword(
		ctx context.Context, usr domain.ChangePasswordParams) error
}

type AuthInteractor struct {
	AuthRepository AuthRepository
	Config         Configurations
	Logger         usecase.Logger
}

func NewAuthInteractor(r AuthRepository, l usecase.Logger) *AuthInteractor {
	return &AuthInteractor{
		AuthRepository: r,
		Config:         *newConfigurations(),
		Logger:         l,
	}
}

type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID   string
	UserRole string
	KeyType  string
	jwt.StandardClaims
}

// Authenticate checks the user credentials
// in request against the db and authenticates the request
func (auth *AuthInteractor) Authenticate(
	reqUser *domain.User, user *domain.User) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(reqUser.Password),
	)
	if err != nil {
		auth.Logger.LogError("password hashes are not same")
		return false
	}
	return true
}

// GenerateAccessToken generates a new access token for the given user
func (auth *AuthInteractor) GenerateAccessToken(
	user *domain.User) (string, error) {

	userID := strconv.FormatInt(user.ID, 10)
	tokenType := "access"
	userRole := user.Role

	claims := AccessTokenCustomClaims{
		userID,
		userRole,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Minute * time.Duration(auth.Config.JwtExpiration),
			).Unix(),
			Issuer: "risc_app.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.Config.AccessTokenPrivateKeyPath)

	if err != nil {
		auth.Logger.LogError("unable to read access private key", err)
		return "", errors.New(
			"could not generate access token. please try again later 1")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.Logger.LogError("unable to parse private key", "error", err)
		return "", errors.New(
			"could not generate access token. please try again later 2")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// ValidateAccessToken parses and validates the given access token
// returns the userId present in the token payload
func (auth *AuthInteractor) ValidateAccessToken(
	tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				auth.Logger.LogError("Unexpected signing method in auth token")
				return nil, errors.New("unexpected signing method in auth token")
			}
			verifyBytes, err := ioutil.ReadFile(auth.Config.AccessTokenPublicKeyPath)
			if err != nil {
				auth.Logger.LogError("unable to read public key", "error", err)
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			if err != nil {
				auth.Logger.LogError("unable to parse public key", "error", err)
				return nil, err
			}

			return verifyKey, nil
		})

	if err != nil {
		auth.Logger.LogError("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, claims.UserRole, nil
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (auth *AuthInteractor) GenerateCustomKey(
	userID string, tokenHash string) string {

	// data := userID + tokenHash
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// GenerateRefreshToken generate a new refresh token for the given user
func (auth *AuthInteractor) GenerateRefreshToken(
	user *domain.User) (string, error) {
	userID := strconv.FormatInt(user.ID, 10)
	cusKey := auth.GenerateCustomKey(userID, user.Tokenhash)
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		userID,
		cusKey,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				24 * time.Hour * time.Duration(auth.Config.JwtRefreshExpiration),
			).Unix(),
			Issuer: "risc_app.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.Config.RefreshTokenPrivateKeyPath)
	if err != nil {
		auth.Logger.LogError("unable to read refresh private key", err)
		return "", errors.New(
			"could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.Logger.LogError("unable to parse private key", "error", err)
		return "", errors.New(
			"could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// ValidateRefreshToken parses and validates the given refresh token
// returns the userId and customkey present in the token payload
func (auth *AuthInteractor) ValidateRefreshToken(
	tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				auth.Logger.LogError("unexpected signing method in auth token")
				return nil, errors.New("unexpected signing method in auth token")
			}
			verifyBytes, err := ioutil.ReadFile(auth.Config.RefreshTokenPublicKeyPath)
			if err != nil {
				auth.Logger.LogError("unable to read public key", "error", err)
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			if err != nil {
				auth.Logger.LogError("unable to parse public key", "error", err)
				return nil, err
			}

			return verifyKey, nil
		})

	if err != nil {
		auth.Logger.LogError("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	auth.Logger.LogAccess("ok", ok)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		auth.Logger.LogAccess("could not extract claims from token")
		return "", "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, claims.CustomKey, nil
}

func (auth *AuthInteractor) RegisterUser(
	ctx context.Context,
	u domain.CreateUserParams) (string, error) {

	hashedPass, err := hashPassword(u.Password)
	if err != nil {
		return "", err
	}

	usr := domain.CreateRegisterUserParams{
		MailVerfyCode:   utils.GenerateRandomString(8),
		MailVerfyExpire: time.Now().Add(720 * time.Hour),
		Name:            u.Name,
		Role:            "user",
		Email:           u.Email,
		Username:        u.Username,
		Password:        hashedPass,
		Tokenhash:       utils.GenerateRandomString(15),
	}

	err = auth.AuthRepository.CreateRegisterUser(ctx, usr)

	return usr.MailVerfyCode, err
}
func hashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func (auth *AuthInteractor) RegisterOauthUser(
	ctx context.Context,
	u domain.CreateOauthUserParams) (domain.User, error) {

	newUsr, err := auth.AuthRepository.CreateOauthUser(ctx, u)

	return newUsr, err
}

func (auth *AuthInteractor) UserByEmail(
	ctx context.Context, email string) (domain.User, error) {

	u, err := auth.AuthRepository.GetUserByEmail(ctx, email)
	return u, err
}

func (auth *AuthInteractor) UserById(
	ctx context.Context, usrID int64) (domain.User, error) {

	u, err := auth.AuthRepository.GetUserById(ctx, usrID)
	return u, err
}

func (auth *AuthInteractor) BasicUserById(
	ctx context.Context, usrID int64) (domain.ShowUserParams, error) {

	u, err := auth.AuthRepository.GetBasicUserById(ctx, usrID)
	return u, err
}

func (auth *AuthInteractor) ShowAllUsers(
	ctx context.Context) ([]domain.ShowUserParams, error) {

	users, err := auth.AuthRepository.GetListusers(ctx)
	return users, err
}

// GetCompleteListusers

func (auth *AuthInteractor) ShowCompleteUsers(
	ctx context.Context) ([]domain.User, error) {

	users, err := auth.AuthRepository.GetCompleteListusers(ctx)
	return users, err
}

func (auth *AuthInteractor) GenerateResetPasswCode(
	ctx context.Context,
	email string) (mail_usecase.Mail, string, error) {

	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	passVerCode := fmt.Sprint(rand.Intn(max-min+1) + min)
	passwordVerfyCode := passVerCode[:3] + "-" + passVerCode[3:]
	passwordVerfyExpire := sql.NullTime{
		Time:  time.Now().Local().Add(1 * time.Hour),
		Valid: true,
	}

	resetPassCodeUpdate := domain.GenerateResetPasswordCodeParams{
		PasswordVerfyCode:   passwordVerfyCode,
		PasswordVerfyExpire: passwordVerfyExpire.Time,
		Email:               email,
	}

	err := auth.AuthRepository.GenerateResetPasswordCode(ctx, resetPassCodeUpdate)

	verfyPassword := mail_usecase.Mail{
		Reciever:  email,
		MailTitle: "Password reset code",
		Type:      2,
	}

	return verfyPassword, passwordVerfyCode, err
}

func (auth *AuthInteractor) UserMailVerify(
	ctx context.Context, email string) error {

	err := auth.AuthRepository.VerifyUserMail(ctx, email)

	return err
}

func (auth *AuthInteractor) UpdatePassword(
	ctx context.Context, usr domain.ChangePasswordParams) error {

	hashNewPass, err := hashPassword(usr.NewPassword)
	usr.NewPassword = hashNewPass
	if err != nil {
		return err
	}

	err = auth.AuthRepository.ChangePassword(ctx, usr)

	return err
}
