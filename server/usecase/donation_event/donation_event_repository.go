package donationevent_usecase

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonationEventRepository interface {
	DonationEventsList(
		ctx context.Context,
		arg domain.DonationEventsListParams) ([]domain.DonationEvent, error)
	DonationEventById(
		ctx context.Context, eventID int64) (domain.DonationEvent, error)
}
