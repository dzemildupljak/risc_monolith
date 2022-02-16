package utils

import (
	"errors"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

type ValidationError struct {
	Message   string
	FieldName string
}
type ValidationErrors struct {
	Fields []ValidationError
}

type AuthValidator struct {
	logger usecase.Logger
}

func NewAuthValidator(l usecase.Logger) *AuthValidator {
	return &AuthValidator{
		logger: l,
	}
}

func (av *AuthValidator) ValidateLoginValues(usr domain.ShowLoginUser) (ValidationErrors, error) {
	valErr := ValidationErrors{Fields: make([]ValidationError, 2)}

	if usr.Password == "" || len(usr.Password) < 5 || usr.Email == "" {
		valErr.Fields[0].Message = "Invalid Email"
		valErr.Fields[0].FieldName = "email"
		valErr.Fields[1].Message = "Invalid Password"
		valErr.Fields[1].FieldName = "password"

		return valErr, errors.New("invalid login values")
	}
	return valErr, nil
}

func (av *AuthValidator) ValidateForgotPassValues(usr domain.ForgotPasswordValues) (ValidationErrors, error) {
	valErr := ValidationErrors{Fields: make([]ValidationError, 2)}

	if usr.New_password == "" ||
		len(usr.New_password) < 5 ||
		usr.Code == "" ||
		usr.New_password_second == "" ||
		usr.New_password != usr.New_password_second {
		valErr.Fields[0].Message = "Invalid Email"
		valErr.Fields[0].FieldName = "email"
		valErr.Fields[1].Message = "Invalid Password"
		valErr.Fields[1].FieldName = "password"

		return valErr, errors.New("invalid login values")
	}
	return valErr, nil
}
