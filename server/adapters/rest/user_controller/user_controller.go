package user_rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/user_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"github.com/gorilla/mux"
)

type UserController struct {
	userInteractor user_usecase.UserInteractor
	logger         usecase.Logger
}

func NewUserController(
	ui user_usecase.UserInteractor,
	logger usecase.Logger) *UserController {

	return &UserController{
		logger:         logger,
		userInteractor: ui,
	}
}

// Index return response which contain a listing of the resource of users.
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.userInteractor.ListUsersInteract(r.Context())

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
func (uc *UserController) UserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["user_id"], 10, 64)
	if err != nil {
		uc.logger.LogError("user Id validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get user. Please try again later",
			})
		return
	}
	usr, err := uc.userInteractor.UserByIdInteract(r.Context(), userId)

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// GetUserByEmail return response of the resource of users.
func (uc *UserController) UserByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	useEmail := params["user_id"]

	usr, err := uc.userInteractor.UserByEmailInteract(r.Context(), useEmail)

	if err != nil {
		uc.logger.LogError("get basic user by email", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get user. Please try again later",
			})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// UpdateUserById update and return response of the resource of users.
func (uc *UserController) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["user_id"], 10, 64)
	if err != nil {
		uc.logger.LogError("user Id validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get user. Please try again later",
			})
		return
	}
	user := &domain.UpdateUserParams{}

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		uc.logger.LogError("deserialization of user json failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create user.Please try again later",
			})
		return
	}
	usr, err := uc.userInteractor.UserUpdate(r.Context(), userId, *user)

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}
