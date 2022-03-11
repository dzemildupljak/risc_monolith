-- name: DonorList :many
SELECT * FROM donors
ORDER BY name;


-- name: DonorByuniqueIdentificationNumber :one
SELECT * FROM donors
WHERE donor_unique_identification_number = $1
LIMIT 1;

-- name: DonorsByBoodType :many
SELECT * FROM donors
WHERE donor_blood_type = $1
LIMIT $2;

-- name: DonorsByBloodTypeNum :many
SELECT * FROM donors
WHERE donor_blood_type_num = $1;


-- name: DonorsWithValidNewDonation :many
SELECT * FROM donors
WHERE donor_last_donation < now() - '4 months' :: interval
ORDER BY donor_surname;