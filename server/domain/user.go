package domain

import (
	"database/sql"
	"time"
)

type User struct {
	ID                  int64        `json:"id"`
	Name                string       `json:"name"`
	Username            string       `json:"username"`
	Email               string       `json:"email" validate:"required"`
	AccessToken         string       `json:"access_token"`
	Password            string       `json:"password" validate:"required"`
	Address             string       `json:"address"`
	Tokenhash           string       `json:"tokenhash"`
	Isverified          bool         `json:"isverified"`
	MailVerfyCode       string       `json:"mail_verfy_code"`
	MailVerfyExpire     sql.NullTime `json:"mail_verfy_expire"`
	PasswordVerfyCode   string       `json:"password_verfy_code"`
	PasswordVerfyExpire sql.NullTime `json:"password_verfy_expire"`
	Createdat           time.Time    `json:"createdat"`
	Updatedat           time.Time    `json:"updatedat"`
}

type CreateUserParams struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Tokenhash string `json:"tokenhash"`
}
type CreateRegisterUserParams struct {
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Tokenhash       string    `json:"tokenhash"`
	MailVerfyCode   string    `json:"mail_verfy_code"`
	MailVerfyExpire time.Time `json:"mail_verfy_expire"`
}

type UpdateUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ShowLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ShowUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type ChangePasswordParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ChangePasswordValues struct {
	Code                string `json:"code"`
	Old_password        string `json:"old_password"`
	New_password        string `json:"new_password"`
	New_password_second string `json:"new_password_second"`
}

type GenerateResetPasswordCodeParams struct {
	PasswordVerfyCode   string    `json:"password_verfy_code"`
	PasswordVerfyExpire time.Time `json:"password_verfy_expire"`
	Email               string    `json:"email"`
}

type ForgotPasswordValues struct {
	Code                string `json:"code"`
	Email               string `json:"email"`
	New_password        string `json:"new_password"`
	New_password_second string `json:"new_password_second"`
}

type UserEmail struct {
	Email string `json:"email"`
}
