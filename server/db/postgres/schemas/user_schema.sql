CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL DEFAULT '',
  username  text NOT NULL DEFAULT '',
  email  text NOT NULL DEFAULT '',
  access_token text NOT NULL DEFAULT '',
  password  text NOT NULL DEFAULT '',
  address text NOT NULL DEFAULT '',
  tokenhash text NOT NULL DEFAULT '',
  isverified bool NOT NULL DEFAULT 'false',
  createdat time NOT NULL DEFAULT CURRENT_TIME,
  updatedat time NOT NULL DEFAULT CURRENT_TIME
);