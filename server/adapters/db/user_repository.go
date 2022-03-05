package psql

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type UserRepository struct {
	Queries Queries
}

func NewUserRepository(q Queries) *UserRepository {
	return &UserRepository{
		Queries: q,
	}

}

const userList = `-- name: UserList :many
SELECT id, name, username, role, email, address, isverified FROM users
ORDER BY name
`

func (q *UserRepository) UserList(ctx context.Context) ([]domain.ShowUserParams, error) {
	rows, err := q.Queries.db.QueryContext(ctx, userList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.ShowUserParams
	for rows.Next() {
		var i domain.ShowUserParams
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.Username,
			&i.Role,
			&i.Email,
			&i.Address,
			&i.Isverified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const completeListusers = `-- name: ListCompleteUsers :many
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, createdat, updatedat FROM users
ORDER BY name
`

func (q *UserRepository) ListCompleteUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := q.Queries.db.QueryContext(ctx, completeListusers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.User
	for rows.Next() {
		var i domain.User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.Role,
			&i.Email,
			&i.AccessToken,
			&i.Password,
			&i.Address,
			&i.Tokenhash,
			&i.Isverified,
			&i.Createdat,
			&i.Updatedat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const userByEmail = `-- name: UserByEmail :one
SELECT id, name, username, role, email, address, isverified FROM users
WHERE email = $1 LIMIT 1
`

func (ac *UserRepository) UserByEmail(ctx context.Context, email string) (domain.ShowUserParams, error) {
	row := ac.Queries.db.QueryRowContext(ctx, userByEmail, email)
	var i domain.ShowUserParams
	err := row.Scan(
		&i.Name,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.Address,
		&i.Isverified,
	)
	return i, err
}

const completeUserById = `-- name: CompleteUserById :one
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, mail_verfy_code, mail_verfy_expire, password_verfy_code, password_verfy_expire, createdat, updatedat FROM users
WHERE id = $1 LIMIT 1
`

func (q *UserRepository) CompleteUserById(ctx context.Context, id int64) (domain.User, error) {
	row := q.Queries.db.QueryRowContext(ctx, completeUserById, id)
	var i domain.User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.AccessToken,
		&i.Password,
		&i.Address,
		&i.Tokenhash,
		&i.Isverified,
		&i.MailVerfyCode,
		&i.MailVerfyExpire,
		&i.PasswordVerfyCode,
		&i.PasswordVerfyExpire,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const userById = `-- name: UserById :one
SELECT id, name, username, role, email, address, isverified FROM users
WHERE id = $1 LIMIT 1
`

func (q *UserRepository) UserById(ctx context.Context, id int64) (domain.ShowUserParams, error) {
	row := q.Queries.db.QueryRowContext(ctx, userById, id)
	var i domain.ShowUserParams
	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.Address,
		&i.Isverified,
	)
	return i, err
}

const userUpdate = `-- name: UpdateUser :one
UPDATE users 
SET name = $2, 
    username = $3, 
    address = $4
WHERE id = $1
RETURNING id, name, username, role, email, address, isverified
`

func (q *UserRepository) UpdateUser(ctx context.Context, usrId int64, arg domain.UpdateUserParams) (domain.ShowUserParams, error) {
	row := q.Queries.db.QueryRowContext(ctx, userUpdate,
		usrId,
		arg.Name,
		arg.Username,
		arg.Address,
	)
	var i domain.ShowUserParams
	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.Address,
		&i.Isverified,
	)

	return i, err
}

const userDeleteById = `-- name: UserDeleteById :exec
DELETE FROM users
WHERE id = $1
`

func (q *UserRepository) UserDeleteById(ctx context.Context, id int64) error {
	_, err := q.Queries.db.ExecContext(ctx, userDeleteById, id)
	return err
}
