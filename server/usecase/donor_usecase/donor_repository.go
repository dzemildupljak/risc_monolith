package donor_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonorRepository interface {
	CreateDonor(ctx context.Context, arg domain.CreateDonorParams) (int64, error)
	DonorList(ctx context.Context) ([]domain.Donor, error)
	DonorById(ctx context.Context, id int64) (domain.Donor, error)
	DonorByuniqueIdentificationNumber(
		ctx context.Context, donorUniqueIdentificationNumber string) (domain.Donor, error)
	DonorsByBloodType(ctx context.Context, donorBloodTypeNum int16) ([]domain.Donor, error)
	DonorsWithValidNewDonation(ctx context.Context) ([]domain.Donor, error)
}
