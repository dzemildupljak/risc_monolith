package psql

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonationEventRepository struct {
	Queries Queries
}

func NewDonationEventRepository(q Queries) *DonationEventRepository {
	return &DonationEventRepository{
		Queries: q,
	}
}

const donationEventsList = `-- name: DonationEventsList :many
SELECT event_id, event_name, event_location, event_start_date, event_end_date, event_organizer, event_organization, event_number_blood_donors FROM donation_events
ORDER BY $1::text
LIMIT $2::integer
`

func (q *Queries) DonationEventsList(
	ctx context.Context,
	arg domain.DonationEventsListParams) ([]domain.DonationEvent, error) {

	rows, err := q.db.QueryContext(ctx, donationEventsList, arg.RowOrder, arg.LimitSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.DonationEvent
	for rows.Next() {
		var i domain.DonationEvent
		if err := rows.Scan(
			&i.EventID,
			&i.EventName,
			&i.EventLocation,
			&i.EventStartDate,
			&i.EventEndDate,
			&i.EventOrganizer,
			&i.EventOrganization,
			&i.EventNumberBloodDonors,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const donationEventById = `-- name: DonationEventById :one
SELECT event_id, event_name, event_location, event_start_date, event_end_date, event_organizer, event_organization, event_number_blood_donors FROM donation_events
WHERE event_id=$1
`

func (q *Queries) DonationEventById(
	ctx context.Context, eventID int64) (domain.DonationEvent, error) {

	row := q.db.QueryRowContext(ctx, donationEventById, eventID)
	var i domain.DonationEvent
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.EventLocation,
		&i.EventStartDate,
		&i.EventEndDate,
		&i.EventOrganizer,
		&i.EventOrganization,
		&i.EventNumberBloodDonors,
	)
	return i, err
}
