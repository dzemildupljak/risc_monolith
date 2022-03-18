package domain

import (
	"database/sql"
	"time"
)

type Donor struct {
	DonorID                         int64        `json:"donor_id"`
	DonorUniqueIdentificationNumber string       `json:"donor_unique_identification_number"`
	DonorName                       string       `json:"donor_name"`
	DonorSurname                    string       `json:"donor_surname"`
	DonorAddress                    string       `json:"donor_address"`
	DonorLastDonation               sql.NullTime `json:"donor_last_donation"`
	DonorPhoneNumber                string       `json:"donor_phone_number"`
	DonorBloodType                  string       `json:"donor_blood_type"`
	DonorBloodTypeNum               int16        `json:"donor_blood_type_num"`
}

type CreateDonorParams struct {
	DonorName                       string       `json:"donor_name"`
	DonorSurname                    string       `json:"donor_surname"`
	DonorBloodType                  string       `json:"donor_blood_type"`
	DonorUniqueIdentificationNumber string       `json:"donor_unique_identification_number"`
	DonorAddress                    string       `json:"donor_address"`
	DonorLastDonation               sql.NullTime `json:"donor_last_donation"`
	DonorPhoneNumber                string       `json:"donor_phone_number"`
	DonorBloodTypeNum               int16        `json:"donor_blood_type_num"`
}
type ShowDonorParams struct {
	DonorName                       string    `json:"donor_name"`
	DonorSurname                    string    `json:"donor_surname"`
	DonorBloodType                  string    `json:"donor_blood_type"`
	DonorUniqueIdentificationNumber string    `json:"donor_unique_identification_number"`
	DonorAddress                    string    `json:"donor_address"`
	DonorLastDonationDate           time.Time `json:"donor_last_donation"`
	DonorPhoneNumber                string    `json:"donor_phone_number"`
	DonorBloodTypeNum               int16     `json:"donor_blood_type_num"`
}

type DonorsByBloodTypeParams struct {
	RowOrder  string `json:"row_order"`
	LimitSize int32  `json:"limit_size"`
}
