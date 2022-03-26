package donor_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type DonorInteractor struct {
	logger          usecase.Logger
	donorRepository DonorRepository
	mapper          utils.Mapper
}

func NewDonorInteractor(
	r DonorRepository,
	l usecase.Logger,
	m utils.Mapper) *DonorInteractor {
	return &DonorInteractor{
		logger:          l,
		donorRepository: r,
		mapper:          m,
	}
}
func (di *DonorInteractor) CreateNewDonor(
	ctx context.Context, donor domain.CreateDonorParamsDTO) error {

	d := domain.CreateDonorParams{
		DonorName:                       donor.DonorName,
		DonorSurname:                    donor.DonorSurname,
		DonorBloodType:                  donor.DonorBloodType,
		DonorUniqueIdentificationNumber: donor.DonorUniqueIdentificationNumber,
		DonorAddress:                    donor.DonorAddress,
		DonorLastDonation: sql.NullTime{
			Valid: true,
			Time:  donor.DonorLastDonation,
		},
		DonorPhoneNumber:  donor.DonorPhoneNumber,
		DonorBloodTypeNum: donor.DonorBloodTypeNum,
	}

	rowNum, err := di.donorRepository.CreateDonor(ctx, d)
	if err != nil || rowNum == 0 {
		return errors.New("donor with given parameters already exists")
	}
	return nil
}

func (di *DonorInteractor) GetAllDonors(
	ctx context.Context) ([]domain.ShowDonorParams, error) {
	donors, err := di.donorRepository.DonorList(ctx)

	var basicDonors []domain.ShowDonorParams

	if err != nil {
		return []domain.ShowDonorParams{}, err
	}

	for _, d := range donors {
		basicDonors = append(basicDonors, di.mapper.BasicDonorMapper(&d))
	}

	return basicDonors, nil
}

func (di *DonorInteractor) GetDonorById(
	ctx context.Context, donor_id int64) (domain.ShowDonorParams, error) {
	donor, err := di.donorRepository.DonorById(ctx, donor_id)

	return di.mapper.BasicDonorMapper(&donor), err
}

func (di *DonorInteractor) GetDonorByIdentificationNum(
	ctx context.Context, id_number string) (domain.ShowDonorParams, error) {
	donor, err := di.donorRepository.DonorByuniqueIdentificationNumber(ctx, id_number)
	return di.mapper.BasicDonorMapper(&donor), err
}

func (di *DonorInteractor) GetDonorsByBloodType(
	ctx context.Context, blood_type int16) ([]domain.ShowDonorParams, error) {
	donors, err := di.donorRepository.DonorsByBloodType(ctx, blood_type)

	var basicDonors []domain.ShowDonorParams

	if err != nil {
		return []domain.ShowDonorParams{}, err
	}

	for _, d := range donors {
		basicDonors = append(basicDonors, di.mapper.BasicDonorMapper(&d))
	}

	return basicDonors, nil
}

func (di *DonorInteractor) GetDonorsWithValidNewDonation(
	ctx context.Context, blood_typ_num int16) ([]domain.ShowDonorParams, error) {
	donors, err := di.donorRepository.DonorsByBloodType(ctx, blood_typ_num)

	var basicDonors []domain.ShowDonorParams

	if err != nil {
		return []domain.ShowDonorParams{}, err
	}

	for _, d := range donors {
		basicDonors = append(basicDonors, di.mapper.BasicDonorMapper(&d))
	}

	return basicDonors, nil

}
