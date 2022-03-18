package auth_rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

// VerificationDataKey is used as the key for storing the VerificationData in context at middleware
type VerificationDataKey struct{}

// VerificationData represents the type for the data stored for verification.
type VerificationData struct {
	Email string `json:"email" validate:"required" sql:"email"`
	Code  string `json:"code" validate:"required" sql:"code"`
	Type  int64  `json:"type" sql:"type"`
}

// MiddlewareValidateUser validates the user in the request
func (ac *AuthController) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		user := &domain.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			ac.logger.LogError("deserialization of user json failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Invalid credentials deserialization",
				})
			return
		}

		err = validEmail(user.Email)

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

		// add the user to the context
		ctx := context.WithValue(r.Context(), auth_usecase.UserKey{}, *user)
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateAccessToken validates whether the request contains a bearer token
// it also decodes and authenticates the given token
func (ac *AuthController) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		ac.logger.LogAccess("validating access token")

		token, err := extractToken(r)
		if err != nil {
			ac.logger.LogError("token not provided or malformed")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		userID, _, err := ac.ai.ValidateAccessToken(token)
		if err != nil {
			ac.logger.LogError("token validation failed1", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}
		ac.logger.LogAccess("access token validated", userID)

		ctx := context.WithValue(r.Context(), auth_usecase.UserIDKey{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateAccessToken validates whether the request contains a bearer token
// it also decodes and authenticates the given token
func (ac *AuthController) MiddlewareValidateAdminAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		ac.logger.LogAccess("validating access token")

		token, err := extractToken(r)
		if err != nil {
			ac.logger.LogError("token not provided or malformed")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		userID, userRole, err := ac.ai.ValidateAccessToken(token)
		if err != nil {
			ac.logger.LogError("token validation failed1", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		if userRole != "admin" {
			ac.logger.LogError("token validation rola failed", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		ac.logger.LogAccess("access token validated", userID)

		ctx := context.WithValue(r.Context(), auth_usecase.UserIDKey{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateRefreshToken validates whether the request contains a bearer token
// it also decodes and authenticates the given token
func (ac *AuthController) MiddlewareValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		ac.logger.LogAccess("validating refresh token")
		token, err := extractToken(r)
		if err != nil {
			ac.logger.LogError("token not provided or malformed")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}
		ac.logger.LogAccess("token present in header")

		userID, customKey, err := ac.ai.ValidateRefreshToken(token)
		if err != nil {
			ac.logger.LogError("token validation failed2", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}
		ac.logger.LogAccess("refresh token validated")

		usrId, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			ac.logger.LogError("token validation failed3", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})
			return
		}
		user, err := ac.ai.UserById(context.Background(), usrId)
		if err != nil {
			ac.logger.LogError("invalid token: wrong userID while parsing", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		actualCustomKey := ac.ai.GenerateCustomKey(
			strconv.FormatInt(user.ID, 10),
			user.Tokenhash)

		if customKey != actualCustomKey {
			ac.logger.LogAccess("wrong token: authetincation failed")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Authentication failed. Token not provided or malformed",
				})

			return
		}

		fmt.Println("middleware refresh token", user)

		ctx := context.WithValue(r.Context(), auth_usecase.UserKey{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (ac *AuthController) MiddlewareValidateMailVerificationData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		ac.logger.LogAccess("validating verification data middleware")

		vals := r.URL.Query()

		c, err := strconv.ParseInt(vals["type"][0], 10, 64)
		if err != nil {
			ac.logger.LogError("deserialization of verification code failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify data try again later"})
			return
		}

		if vals["email"][0] == "" || vals["code"][0] == "" {
			ac.logger.LogError("deserialization of verification data failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify data try again later"})
			return
		}

		verificationData := &VerificationData{
			Email: vals["email"][0],
			Code:  vals["code"][0],
			Type:  c,
		}

		user, err := ac.ai.UserByEmail(r.Context(), verificationData.Email)
		if err != nil {
			ac.logger.LogError("verification code failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify data try again later"})
			return
		}

		if verificationData.Type != 1 {
			ac.logger.LogError("verification code failed1", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify mail try again later"})
			return
		}

		validateExTime := utils.ValidateExpirationTime(user.MailVerfyExpire.Time)

		if !validateExTime {
			ac.logger.LogError("verification code failed2")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify data try again later"})
			return
		}

		if user.MailVerfyCode != verificationData.Code {
			ac.logger.LogError("verification code failed3", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Failed to verify data try again later"})
			return
		}

		fmt.Println("end middleware")

		// add the ValidationData to context
		ctx := context.WithValue(r.Context(), VerificationDataKey{}, *verificationData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("token not provided or malformed")
	}

	return authHeaderContent[1], nil
}

func validEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
