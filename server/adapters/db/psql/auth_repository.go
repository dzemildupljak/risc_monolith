package psql

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/lib/pq"
)

type AuthRepository struct {
	Queries Queries
}

func NewAuthRepository(q Queries) *AuthRepository {
	return &AuthRepository{
		Queries: q,
	}

}

const createRegisterUser = `-- name: CreateRegisterUser :one
INSERT INTO users (
  name, email, username, role, password, tokenhash, mail_verfy_code, mail_verfy_expire, updatedat, createdat
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING id, name, username, role, email, access_token, password, address, tokenhash, isverified, mail_verfy_code, mail_verfy_expire, createdat, updatedat
`

func (q *AuthRepository) CreateRegisterUser(ctx context.Context, arg domain.CreateRegisterUserParams) error {
	row := q.Queries.db.QueryRowContext(ctx, createRegisterUser,
		arg.Name,
		arg.Email,
		arg.Username,
		arg.Role,
		arg.Password,
		arg.Tokenhash,
		arg.MailVerfyCode,
		arg.MailVerfyExpire,
	)
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
		&i.Createdat,
		&i.Updatedat,
	)
	return err
}

const createOauthUser = `-- name: CreateOauthUser :one
INSERT INTO users (
	email, tokenhash, isverified, oauth_id, updatedat, createdat
) VALUES (
	$1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING id, name, username, role, email, access_token, password, address, tokenhash, isverified, oauth_id, mail_verfy_code, mail_verfy_expire, password_verfy_code, password_verfy_expire, createdat, updatedat`

func (q *AuthRepository) CreateOauthUser(ctx context.Context, arg domain.CreateOauthUserParams) (domain.User, error) {
	row := q.Queries.db.QueryRowContext(ctx,
		createOauthUser,
		arg.Email,
		arg.Tokenhash,
		arg.Isverified,
		pq.Array(arg.OauthID))

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
		pq.Array(&i.OauthID),
		&i.MailVerfyCode,
		&i.MailVerfyExpire,
		&i.PasswordVerfyCode,
		&i.PasswordVerfyExpire,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1
`

func (q *AuthRepository) DeleteUserById(ctx context.Context, id int64) error {
	_, err := q.Queries.db.ExecContext(ctx, deleteUserById, id)
	return err
}

const getListusers = `-- name: GetListusers :many
SELECT id, name, username, role, email, address, isverified FROM users
ORDER BY name
`

func (q *AuthRepository) GetListusers(ctx context.Context) ([]domain.ShowUserParams, error) {
	rows, err := q.Queries.db.QueryContext(ctx, getListusers)
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

const getCompleteListusers = `-- name: GetListusers :many
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, createdat, updatedat FROM users
ORDER BY name
`

func (q *AuthRepository) GetCompleteListusers(ctx context.Context) ([]domain.User, error) {
	rows, err := q.Queries.db.QueryContext(ctx, getCompleteListusers)
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

const getOneUser = `-- name: GetOneUser :one
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, createdat, updatedat FROM users
WHERE id = $1 LIMIT 1
`

func (q *AuthRepository) GetOneUser(ctx context.Context, id int64) (domain.User, error) {
	row := q.Queries.db.QueryRowContext(ctx, getOneUser, id)
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
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, username, role, email,access_token, password, address, tokenhash, isverified, mail_verfy_code, mail_verfy_expire, password_verfy_code, password_verfy_expire, createdat, updatedat FROM users
WHERE email = $1 LIMIT 1
`

func (ac *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row := ac.Queries.db.QueryRowContext(ctx, getUserByEmail, email)
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

const getLogedUserByEmaiAuth = `-- name: GetLogedUserByEmai :one
SELECT email,password FROM users
WHERE email = $1 LIMIT 1
`

func (q *AuthRepository) GetLogedUserByEmai(ctx context.Context, username string) (domain.ShowLoginUser, error) {
	row := q.Queries.db.QueryRowContext(ctx, getLogedUserByEmaiAuth, username)
	var i domain.ShowLoginUser
	err := row.Scan(&i.Email, &i.Password)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET name = $1, 
    username = $2, 
    address = $3
RETURNING id, name, username, role, email, access_token, password, address, tokenhash, isverified, createdat, updatedat
`

func (q *AuthRepository) UpdateUser(ctx context.Context, arg domain.UpdateUserParams) (domain.User, error) {
	row := q.Queries.db.QueryRowContext(ctx, updateUser,
		arg.Name,
		arg.Username,
		arg.Address,
	)
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
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, name, username, role, email, access_token, password, address, tokenhash, isverified, mail_verfy_code, mail_verfy_expire, password_verfy_code, password_verfy_expire, createdat, updatedat FROM users
WHERE id = $1 LIMIT 1
`

func (q *AuthRepository) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	row := q.Queries.db.QueryRowContext(ctx, getUserById, id)
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

const getBasicUserById = `-- name: GetUserById :one
SELECT id, name, username, role, email, address, isverified FROM users
WHERE id = $1 LIMIT 1
`

func (q *AuthRepository) GetBasicUserById(ctx context.Context, id int64) (domain.ShowUserParams, error) {
	row := q.Queries.db.QueryRowContext(ctx, getBasicUserById, id)
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

const verifyUserMail = `-- name: VerifyUserMail :exec
UPDATE users
SET isverified = true
WHERE email = $1
`

func (q *AuthRepository) VerifyUserMail(ctx context.Context, email string) error {
	_, err := q.Queries.db.ExecContext(ctx, verifyUserMail, email)
	return err
}

const generateResetPasswordCode = `-- name: GenerateResetPasswordCode :exec
UPDATE users
SET password_verfy_code = $1, password_verfy_expire = $2
WHERE email = $3
`

func (q *AuthRepository) GenerateResetPasswordCode(ctx context.Context, arg domain.GenerateResetPasswordCodeParams) error {
	_, err := q.Queries.db.ExecContext(ctx, generateResetPasswordCode, arg.PasswordVerfyCode, arg.PasswordVerfyExpire, arg.Email)
	return err
}

const changePassword = `-- name: ChangePassword :exec
UPDATE users
SET password = $1
WHERE email = $2
`

func (q *AuthRepository) ChangePassword(ctx context.Context, arg domain.ChangePasswordParams) error {
	_, err := q.Queries.db.ExecContext(ctx, changePassword, arg.NewPassword, arg.Email)
	return err
}
