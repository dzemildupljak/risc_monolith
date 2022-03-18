package donor_usecase

import (
	"context"
	"errors"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

type DonorInteractor struct {
	logger          usecase.Logger
	donorRepository DonorRepository
}

func NewDonorInteractor(r DonorRepository, l usecase.Logger) *DonorInteractor {
	return &DonorInteractor{
		logger:          l,
		donorRepository: r,
	}
}
func (di *DonorInteractor) CreateNewDonor(
	ctx context.Context, donor domain.CreateDonorParams) error {

	rowNum, err := di.donorRepository.CreateDonor(ctx, donor)
	if err != nil || rowNum == 0 {
		return errors.New("donor with given parameters exists")
	}
	return nil
}

func (di *DonorInteractor) GetAllDonors(
	ctx context.Context) ([]domain.ShowDonorParams, error) {
	donors, err := di.donorRepository.DonorList(ctx)

	var dnrs []domain.ShowDonorParams

	if err != nil {
		return []domain.ShowDonorParams{}, err
	}

	for _, d := range donors {
		dnrs = append(dnrs, domain.ShowDonorParams{
			DonorName:                       d.DonorName,
			DonorSurname:                    d.DonorSurname,
			DonorBloodType:                  d.DonorBloodType,
			DonorUniqueIdentificationNumber: d.DonorUniqueIdentificationNumber,
			DonorAddress:                    d.DonorAddress,
			DonorPhoneNumber:                d.DonorPhoneNumber,
			DonorBloodTypeNum:               d.DonorBloodTypeNum,
			DonorLastDonationDate:           d.DonorLastDonation.Time,
		})
	}

	return dnrs, nil
}

func (di *DonorInteractor) GetDonorById(
	ctx context.Context, donor_id int64) (domain.Donor, error) {

	return di.donorRepository.DonorById(ctx, donor_id)
}

func (di *DonorInteractor) GetDonorByIdentificationNum(
	ctx context.Context, id_number string) (domain.Donor, error) {

	return di.donorRepository.DonorByuniqueIdentificationNumber(ctx, id_number)
}

func (di *DonorInteractor) GetDonorsByBloodType(
	ctx context.Context, blood_type int16) ([]domain.Donor, error) {

	return di.donorRepository.DonorsByBloodType(ctx, blood_type)
}

func (di *DonorInteractor) GetDonorsWithValidNewDonation(
	ctx context.Context) ([]domain.Donor, error) {

	return di.donorRepository.DonorsWithValidNewDonation(ctx)
}