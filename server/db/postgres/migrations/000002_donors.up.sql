CREATE TABLE donors (
  donor_id BIGSERIAL PRIMARY KEY,
  donor_unique_identification_number text NOT NULL DEFAULT '' UNIQUE,
  donor_name text NOT NULL DEFAULT '',
  donor_surname text NOT NULL DEFAULT '',
  donor_address text NOT NULL DEFAULT '',
  donor_last_donation date,
  donor_phone_number text NOT NULL DEFAULT '',
  donor_blood_type text NOT NULL DEFAULT '',
  donor_blood_type_num smallint
);