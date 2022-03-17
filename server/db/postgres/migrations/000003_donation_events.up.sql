CREATE TABLE donation_events (
  event_id BIGSERIAL PRIMARY KEY,
  event_name text NOT NULL DEFAULT '',
  event_location text NOT NULL DEFAULT '',
  event_start_date date,
  event_end_date date,
  event_organizer text NOT NULL DEFAULT '',
  event_organization text NOT NULL DEFAULT '',
  event_number_blood_donors integer,
  createdat time NOT NULL DEFAULT CURRENT_TIME,
  updatedat time NOT NULL DEFAULT CURRENT_TIME
)