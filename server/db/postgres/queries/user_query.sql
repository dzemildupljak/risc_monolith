-- name: GetOneUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

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


-- name: UpdateUser :one
UPDATE users 
SET name = $1, 
    username = $2, 
    email = $3, 
    password = $4
RETURNING *;



-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;