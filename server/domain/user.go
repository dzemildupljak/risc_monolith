package domain

import (
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	AccessToken     string    `json:"access_token"`
	Password        string    `json:"password"`
	Address         string    `json:"address"`
	Tokenhash       string    `json:"tokenhash"`
	Isverified      bool      `json:"isverified"`
	MailVerfyCode   string    `json:"mail_verfy_code"`
	MailVerfyExpire time.Time `json:"mail_verfy_expire"`
	Createdat       time.Time `json:"createdat"`
	Updatedat       time.Time `json:"updatedat"`
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
	Username string `json:"username"`
	Password string `json:"password"`
}

type ShowUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}
