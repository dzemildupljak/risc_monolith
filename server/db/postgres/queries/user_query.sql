-- name: UserList :many
SELECT id, name, username, role, email, address, isverified FROM users
ORDER BY name

-- name: ListCompleteUsers :many
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, createdat, updatedat FROM users
ORDER BY name

-- name: UserByEmail :one
SELECT id, name, username, role, email, address, isverified FROM users
WHERE email = $1 LIMIT 1

-- name: CompleteUserById :one
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, mail_verfy_code, mail_verfy_expire, password_verfy_code, password_verfy_expire, createdat, updatedat FROM users
WHERE id = $1 LIMIT 1

-- name: UserById :one
SELECT id, name, username, role, email, address, isverified FROM users
WHERE id = $1 LIMIT 1

-- name: UpdateUser :one
UPDATE users 
SET name = $2, 
    username = $3, 
    address = $4
WHERE id = $1
RETURNING id, name, username, role, email, address, isverified

-- name: UserDeleteById :exec
DELETE FROM users
WHERE id = $1