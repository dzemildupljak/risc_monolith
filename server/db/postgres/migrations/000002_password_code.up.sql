ALTER TABLE users
ADD COLUMN password_verfy_code text NOT NULL DEFAULT '',
ADD COLUMN password_verfy_expire timestamp with time zone,;