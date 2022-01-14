package auth_rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase/auth_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

// MiddlewareValidateUser validates the user in the request
func (ac *AuthController) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		user := &domain.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			ac.logger.LogError("deserialization of user json failed", "error", err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
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
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}

		userID, err := ac.authInteractor.ValidateAccessToken(token)
		if err != nil {
			ac.logger.LogError("token validation failed1", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}
		ac.logger.LogAccess("access token validated")

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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}
		ac.logger.LogAccess("token present in header")

		userID, customKey, err := ac.authInteractor.ValidateRefreshToken(token)
		if err != nil {
			ac.logger.LogError("token validation failed2", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}
		ac.logger.LogAccess("refresh token validated")

		usrId, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			ac.logger.LogError("token validation failed3", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}
		user, err := ac.authInteractor.UserById(context.Background(), usrId)
		if err != nil {
			ac.logger.LogError("invalid token: wrong userID while parsing", err)
			w.WriteHeader(http.StatusBadRequest)
			// data.ToJSON(&GenericError{Error: "invalid token: authentication failed"}, w)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}

		actualCustomKey := ac.authInteractor.GenerateCustomKey(strconv.FormatInt(user.ID, 10), user.Tokenhash)
		if customKey != actualCustomKey {
			ac.logger.LogAccess("wrong token: authetincation failed")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&utils.GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"})

			return
		}

		ctx := context.WithValue(r.Context(), auth_usecase.UserKey{}, user)
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
