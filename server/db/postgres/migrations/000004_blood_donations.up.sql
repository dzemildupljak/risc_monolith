CREATE TABLE blood_donations (
  donation_id BIGSERIAL PRIMARY KEY,
  donor_id BIGSERIAL  NOT NULL,
  FOREIGN KEY(donor_id) REFERENCES donors(donor_id),
  donation_time timestamp with time zone,
  donation_date date NOT NULL,
  medical_worker text NOT NULL DEFAULT '',
  medical_worker_id BIGSERIAL  NOT NULL,
  FOREIGN KEY(medical_worker_id) REFERENCES users(id),
  donation_event_id BIGSERIAL  NOT NULL,
  FOREIGN KEY(donation_event_id) REFERENCES donation_events(event_id)
);
