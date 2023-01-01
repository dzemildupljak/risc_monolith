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
	OauthID             []string     `json:"oauth_id"`
	Role                string       `json:"role"`
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
	Role            string    `json:"role"`
	Password        string    `json:"password"`
	Tokenhash       string    `json:"tokenhash"`
	MailVerfyCode   string    `json:"mail_verfy_code"`
	MailVerfyExpire time.Time `json:"mail_verfy_expire"`
}

type CreateOauthUserParams struct {
	Email      string   `json:"email"`
	Role       string   `json:"role"`
	Tokenhash  string   `json:"tokenhash"`
	Isverified bool     `json:"isverified"`
	OauthID    []string `json:"oauth_id"`
}

type UpdateUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Address  string `json:"address"`
}

type ShowLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ShowUserParams struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Isverified bool   `json:"isverified"`
}

type ChangePasswordParams struct {
	NewPassword string `json:"new_password"`
	Email       string `json:"email"`
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

type UserRole struct {
	Role string `json:"role"`
}
