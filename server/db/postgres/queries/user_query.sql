-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1

-- name: GetUserByUsername :one
SELECT username, password FROM users
WHERE username = $1 LIMIT 1

-- name: GetListusers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  name, username, email, password
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: CreateRegisterUser :one
INSERT INTO users (
  name, email, username, password, tokenhash, updatedat, createdat
) VALUES (
  $1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET name = $1, 
    username = $2, 
    email = $3, 
    password = $4
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1


-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;