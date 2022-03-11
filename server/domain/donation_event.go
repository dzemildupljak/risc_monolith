package domain

import "database/sql"

type DonationEvent struct {
	EventID                int64         `json:"event_id"`
	EventName              string        `json:"event_name"`
	EventLocation          string        `json:"event_location"`
	EventStartDate         sql.NullTime  `json:"event_start_date"`
	EventEndDate           sql.NullTime  `json:"event_end_date"`
	EventOrganizer         string        `json:"event_organizer"`
	EventOrganization      string        `json:"event_organization"`
	EventNumberBloodDonors sql.NullInt32 `json:"event_number_blood_donors"`
}

type DonationEventsListParams struct {
	RowOrder  string `json:"row_order"`
	LimitSize int32  `json:"limit_size"`
}
