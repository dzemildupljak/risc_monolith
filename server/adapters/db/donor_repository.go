package psql

import (
	"context"

	"github.com/dzemildupljak/risc_monolith/server/domain"
)

type DonorRepository struct {
	Queries Queries
}

func NewDonorRepository(q Queries) *DonorRepository {
	return &DonorRepository{
		Queries: q,
	}
}

const donorList = `-- name: DonorList :many
SELECT donor_id, donor_unique_identification_number, donor_name, donor_surname, donor_address, donor_last_donation, donor_phone_number, donor_blood_type, donor_blood_type_num FROM donors
ORDER BY name
`

func (q *DonorRepository) DonorList(ctx context.Context) ([]domain.Donor, error) {
	rows, err := q.Queries.db.QueryContext(ctx, donorList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorID,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodType,
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

const donorByuniqueIdentificationNumber = `-- name: DonorByuniqueIdentificationNumber :one
SELECT donor_id, donor_unique_identification_number, donor_name, donor_surname, donor_address, donor_last_donation, donor_phone_number, donor_blood_type, donor_blood_type_num FROM donors
WHERE donor_unique_identification_number = $1
LIMIT 1
`

func (q *DonorRepository) DonorByuniqueIdentificationNumber(
	ctx context.Context, donorUniqueIdentificationNumber string) (domain.Donor, error) {

	row := q.Queries.db.QueryRowContext(ctx, donorByuniqueIdentificationNumber, donorUniqueIdentificationNumber)
	var i domain.Donor
	err := row.Scan(
		&i.DonorID,
		&i.DonorUniqueIdentificationNumber,
		&i.DonorName,
		&i.DonorSurname,
		&i.DonorAddress,
		&i.DonorLastDonation,
		&i.DonorPhoneNumber,
		&i.DonorBloodType,
		&i.DonorBloodTypeNum,
	)
	return i, err
}

const donorsByBloodType = `-- name: DonorsByBloodType :many
SELECT donor_id, donor_unique_identification_number, donor_name, donor_surname, donor_address, donor_last_donation, donor_phone_number, donor_blood_type, donor_blood_type_num FROM donors
WHERE donor_blood_type = $1::text
LIMIT $2::integer
`

func (q *Queries) DonorsByBloodType(
	ctx context.Context, arg domain.DonorsByBloodTypeParams) ([]domain.Donor, error) {

	rows, err := q.db.QueryContext(ctx, donorsByBloodType, arg.RowOrder, arg.LimitSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorID,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodType,
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

const donorsByBloodTypeNum = `-- name: DonorsByBloodTypeNum :many
SELECT donor_id, donor_unique_identification_number, donor_name, donor_surname, donor_address, donor_last_donation, donor_phone_number, donor_blood_type, donor_blood_type_num FROM donors
WHERE donor_blood_type_num = $1
`

func (q *DonorRepository) DonorsByBloodTypeNum(
	ctx context.Context, donorBloodTypeNum int16) ([]domain.Donor, error) {

	rows, err := q.Queries.db.QueryContext(ctx, donorsByBloodTypeNum, donorBloodTypeNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorID,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodType,
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

const donorsWithValidNewDonation = `-- name: DonorsWithValidNewDonation :many
SELECT donor_id, donor_unique_identification_number, donor_name, donor_surname, donor_address, donor_last_donation, donor_phone_number, donor_blood_type, donor_blood_type_num FROM donors
WHERE donor_last_donation < now() - '4 months' :: interval
ORDER BY donor_surname
`

func (q *DonorRepository) DonorsWithValidNewDonation(
	ctx context.Context) ([]domain.Donor, error) {

	rows, err := q.Queries.db.QueryContext(ctx, donorsWithValidNewDonation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Donor
	for rows.Next() {
		var i domain.Donor
		if err := rows.Scan(
			&i.DonorID,
			&i.DonorUniqueIdentificationNumber,
			&i.DonorName,
			&i.DonorSurname,
			&i.DonorAddress,
			&i.DonorLastDonation,
			&i.DonorPhoneNumber,
			&i.DonorBloodType,
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
