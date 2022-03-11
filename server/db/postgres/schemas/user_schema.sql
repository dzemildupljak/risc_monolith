CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text NOT NULL DEFAULT '',
  username  text NOT NULL DEFAULT '',
  email  text NOT NULL DEFAULT '' UNIQUE,
  access_token text NOT NULL DEFAULT '',
  password  text NOT NULL DEFAULT '',
  address text NOT NULL DEFAULT '',
  tokenhash text NOT NULL DEFAULT '',
  isverified bool NOT NULL DEFAULT 'false',
  oauth_id text[],
  role text NOT NULL DEFAULT '',
  mail_verfy_code text NOT NULL DEFAULT '',
  mail_verfy_expire timestamp with time zone,
  password_verfy_code text NOT NULL DEFAULT '',
  password_verfy_expire timestamp with time zone,
  createdat time NOT NULL DEFAULT CURRENT_TIME,
  updatedat time NOT NULL DEFAULT CURRENT_TIME
);