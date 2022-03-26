package donor_rest

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonorUsecase interface {
	CreateNewDonor(
		ctx context.Context, donor domain.CreateDonorParamsDTO) error
	GetAllDonors(
		ctx context.Context) ([]domain.ShowDonorParams, error)
	GetDonorById(
		ctx context.Context, donor_id int64) (domain.ShowDonorParams, error)
	GetDonorByIdentificationNum(
		ctx context.Context, id_number string) (domain.ShowDonorParams, error)
	GetDonorsByBloodType(
		ctx context.Context, blood_type int16) ([]domain.ShowDonorParams, error)
	GetDonorsWithValidNewDonation(
		ctx context.Context, blood_typ_num int16) ([]domain.ShowDonorParams, error)
}
