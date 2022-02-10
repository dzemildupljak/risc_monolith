CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL DEFAULT '',
  username  text NOT NULL DEFAULT '',
  email  text NOT NULL DEFAULT '' UNIQUE,
  access_token text NOT NULL DEFAULT '',
  password  text NOT NULL DEFAULT '',
  address text NOT NULL DEFAULT '',
  tokenhash text NOT NULL DEFAULT '',
  isverified bool NOT NULL DEFAULT 'false',
  mail_verfy_code text NOT NULL DEFAULT '',
  mail_verfy_expire time,
  createdat time NOT NULL DEFAULT CURRENT_TIME,
  updatedat time NOT NULL DEFAULT CURRENT_TIME
);