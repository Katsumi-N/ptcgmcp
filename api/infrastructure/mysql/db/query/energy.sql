-- name: EnergyFindById :one
SELECT *
FROM energies
WHERE id = ?;
