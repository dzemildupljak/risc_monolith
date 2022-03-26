package utils

import "github.com/dzemildupljak/risc_monolith/server/domain"

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) BasicDonorMapper(donor *domain.Donor) domain.ShowDonorParams {
	return domain.ShowDonorParams{
		DonorName:                       donor.DonorName,
		DonorSurname:                    donor.DonorSurname,
		DonorBloodType:                  donor.DonorBloodType,
		DonorUniqueIdentificationNumber: donor.DonorUniqueIdentificationNumber,
		DonorAddress:                    donor.DonorAddress,
		DonorPhoneNumber:                donor.DonorPhoneNumber,
		DonorBloodTypeNum:               donor.DonorBloodTypeNum,
		DonorLastDonationDate:           donor.DonorLastDonation.Time,
	}
}
