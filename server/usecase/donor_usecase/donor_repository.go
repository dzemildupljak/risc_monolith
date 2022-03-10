package donor_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonorRepository interface {
	DonorList(ctx context.Context) ([]domain.Donor, error)
	DonorById(ctx context.Context, id int64) (domain.Donor, error)
	DonorByuniqueIdentificationNumber(
		ctx context.Context, donorUniqueIdentificationNumber string) (domain.Donor, error)
	DonorByBoodType(ctx context.Context, donorBloodType string) (domain.Donor, error)
	DonorsWithValidNewDonation(ctx context.Context) ([]domain.Donor, error)
	DonorsByBloodTypeNum(ctx context.Context, donorBloodTypeNum int16) ([]domain.Donor, error)
}
