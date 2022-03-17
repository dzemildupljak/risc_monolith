package psql

import (
	"context"
	"os"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

var encryp_key string

type DonorRepository struct {
	Queries Queries
}

var completeDonor string

func NewDonorRepository(q Queries) *DonorRepository {
	encryp_key = os.Getenv("DB_ENCRYPT_KEY")
	completeDonor = `
	pgp_sym_decrypt(donor_name::bytea, '` + encryp_key + `') as donor_name,
	pgp_sym_decrypt(donor_surname::bytea, '` + encryp_key + `') as donor_surname,
	pgp_sym_decrypt(donor_blood_type::bytea, '` + encryp_key + `') as donor_blood_type,
	pgp_sym_decrypt(donor_unique_identification_number::bytea, '` + encryp_key + `') as donor_unique_identification_number,
	donor_address,
	donor_last_donation,
	donor_phone_number,
	donor_blood_type_num`

	return &DonorRepository{
		Queries: q,
	}
}

func (q *DonorRepository) CreateDonor(ctx context.Context, arg domain.CreateDonorParams) (int64, error) {
	var createDonor string = `-- name: CreateDonor :exec
		INSERT INTO donors(
				donor_name, 
				donor_surname, 
				donor_blood_type,
				donor_unique_identification_number,
				donor_address,
				donor_last_donation,
				donor_phone_number,
				donor_blood_type_num)
		SELECT pgp_sym_encrypt($1, '` + encryp_key + `'), 
				pgp_sym_encrypt($2, '` + encryp_key + `'),
				pgp_sym_encrypt($3, '` + encryp_key + `'),
				pgp_sym_encrypt($4, '` + encryp_key + `'),
				$5, $6, $7, $8
		WHERE NOT EXISTS (
			SELECT donor_id, donor_unique_identification_number FROM donors
			WHERE pgp_sym_decrypt(donor_unique_identification_number::bytea, '` + encryp_key + `') = $4
		)`

	result, err := q.Queries.db.ExecContext(ctx, createDonor,
		arg.DonorName,
		arg.DonorSurname,
		arg.DonorBloodType,
		arg.DonorUniqueIdentificationNumber,
		arg.DonorAddress,
		arg.DonorLastDonation,
		arg.DonorPhoneNumber,
		arg.DonorBloodTypeNum,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (q *DonorRepository) DonorList(ctx context.Context) ([]domain.Donor, error) {
	var donorList string = `-- name: DonorList :many
	SELECT ` + completeDonor + ` 
	FROM donors
	ORDER BY donor_surname
	`

	rows, err := q.Queries.db.QueryContext(ctx, donorList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorBloodType,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodTypeNum,
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

func (q *Queries) DonorById(ctx context.Context, donorID int64) (domain.Donor, error) {
	var donorById string = `-- name: DonorById :one
	SELECT ` + completeDonor + `
	FROM donors
	WHERE donor_id = $1
	`
	row := q.db.QueryRowContext(ctx, donorById, donorID)
	var i domain.Donor
	err := row.Scan(
		&i.DonorName,
		&i.DonorSurname,
		&i.DonorBloodType,
		&i.DonorUniqueIdentificationNumber,
		&i.DonorAddress,
		&i.DonorLastDonation,
		&i.DonorPhoneNumber,
		&i.DonorBloodTypeNum,
	)
	return i, err
}

func (q *Queries) DonorByuniqueIdentificationNumber(ctx context.Context, donorUniqueIdentificationNumber string) (domain.Donor, error) {
	var donorByuniqueIdentificationNumber string = `-- name: DonorByuniqueIdentificationNumber :one
	SELECT ` + completeDonor + `
	FROM donors
	WHERE  pgp_sym_decrypt(donor_unique_identification_number::bytea, '` + encryp_key + `') = $1
	`
	row := q.db.QueryRowContext(ctx, donorByuniqueIdentificationNumber, donorUniqueIdentificationNumber)
	var i domain.Donor
	err := row.Scan(
		&i.DonorName,
		&i.DonorSurname,
		&i.DonorBloodType,
		&i.DonorUniqueIdentificationNumber,
		&i.DonorAddress,
		&i.DonorLastDonation,
		&i.DonorPhoneNumber,
		&i.DonorBloodTypeNum,
	)
	return i, err
}

func (q *Queries) DonorsByBloodType(ctx context.Context, donorBloodTypeNum int16) ([]domain.Donor, error) {
	var donorsByBloodType string = `-- name: DonorsByBloodType :many
	SELECT ` + completeDonor + `
	FROM donors
	WHERE  donor_blood_type_num = $1
	`
	rows, err := q.db.QueryContext(ctx, donorsByBloodType, donorBloodTypeNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorBloodType,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodTypeNum,
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

func (q *Queries) DonorsWithValidNewDonation(ctx context.Context) ([]domain.Donor, error) {
	var donorsWithValidNewDonation string = `-- name: DonorsWithValidNewDonation :many
	SELECT ` + completeDonor + `
	FROM donors
	WHERE donor_last_donation < now() - '4 months' :: interval
	ORDER BY donor_surname
	`
	rows, err := q.db.QueryContext(ctx, donorsWithValidNewDonation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorBloodType,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodTypeNum,
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
