package utils

import (
	"database/sql"
	"time"
)

// This text will appear as description of your response body.
// swagger:response genericResponse
type GenericResponseWrapper struct {
	// in:body
	Body GenericResponse
}
type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AuthResponse
// swagger:response authResponse
type AuthResponseWrapper struct {
	// in:body
	Body AuthResponse
}
type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Email        string `json:"email"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// User with basic info
// swagger:response userResponse
type UserResponseWrapper struct {
	// in:body
	Body struct {
		User showUsers
	}
}

// List of users with basic info
// swagger:response userListResponse
type ListUserResponseWrapper struct {
	// in:body
	Body struct {
		Users []showUsers
	}
}

// List of donors with basic info
// swagger:response donorListResponse
type ListDonorResponseWrapper struct {
	// in:body
	Body struct {
		Donors []showDonor
	}
}

type showDonor struct {
	DonorName                       string    `json:"donor_name"`
	DonorSurname                    string    `json:"donor_surname"`
	DonorBloodType                  string    `json:"donor_blood_type"`
	DonorUniqueIdentificationNumber string    `json:"donor_unique_identification_number"`
	DonorAddress                    string    `json:"donor_address"`
	DonorLastDonationDate           time.Time `json:"donor_last_donation"`
	DonorPhoneNumber                string    `json:"donor_phone_number"`
	DonorBloodTypeNum               int16     `json:"donor_blood_type_num"`
}

type showUsers struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Isverified bool   `json:"isverified"`
}

// List of users with complete info
// swagger:response userListCompleteResponse
type ListCompleteUserResponseWrapper struct {
	// in:body
	Body struct {
		Users []user
	}
}

type user struct {
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
