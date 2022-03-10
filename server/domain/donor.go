package domain

import (
	"database/sql"
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
