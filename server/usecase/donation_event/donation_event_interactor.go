package donationevent_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

type DonationEventUsecase interface {
	GetDonationEventsList(
		ctx context.Context,
		arg domain.DonationEventsListParams) ([]domain.DonationEvent, error)
	GetDonationEventById(
		ctx context.Context, eventID int64) (domain.DonationEvent, error)
}

type DonationEventInteractor struct {
	logger       usecase.Logger
	donationRepo DonationEventRepository
}

func NewDonationEventInteractor(
	d DonationEventRepository,
	l usecase.Logger) *DonationEventInteractor {
	return &DonationEventInteractor{
		logger:       l,
		donationRepo: d,
	}
}

func (dei *DonationEventInteractor) GetDonationEventsList(
	ctx context.Context,
	param domain.DonationEventsListParams) (
	[]domain.DonationEvent, error) {

	return dei.donationRepo.DonationEventsList(ctx, param)
}

func (dei *DonationEventInteractor) GetDonationEventById(
	ctx context.Context,
	donationId int64) (
	domain.DonationEvent, error) {

	return dei.donationRepo.DonationEventById(ctx, donationId)
}
