package docs

import "time"

// swagger:route GET /donor/donors Donor ListDonors
// You receive an list of donors
// responses:
//   200: donorListResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:route GET /donor/{donor_id} Donor DonorById
// You receive an donor
// responses:
//   200: donorResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:route GET /donor/blood-type/{blood_type} Donor DonorByBloodType
// You receive an donor for given id
// responses:
//   200: donorResponse
//   500: genericResponse

/////////////////////////////////////////////////
// swagger:parameters DonorByUniqueNumber
type donorUniqueNumRequest struct {
	// in:body
	Body struct {
		UniqueNumber string `json:"unique_number"`
	}
}

// swagger:route POST /donor/unique Donor DonorByUniqueNumber
// You receive an donor for given Unique Number
// responses:
//   200: donorResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:parameters CreateNewDonor
type donorNewDonorRequest struct {
	// in:body
	Body struct {
		DonorName                       string    `json:"donor_name"`
		DonorSurname                    string    `json:"donor_surname"`
		DonorBloodType                  string    `json:"donor_blood_type"`
		DonorUniqueIdentificationNumber string    `json:"donor_unique_identification_number"`
		DonorAddress                    string    `json:"donor_address"`
		DonorLastDonationDate           time.Time `json:"donor_last_donation"`
		DonorPhoneNumber                string    `json:"donor_phone_number"`
		DonorBloodTypeNum               int16     `json:"donor_blood_type_num"`
	}
}

// swagger:route POST /donor/ Donor CreateNewDonor
// You create an donor for given values
// responses:
//   200: genericResponse
//   500: genericResponse

/////////////////////////////////////////////////
