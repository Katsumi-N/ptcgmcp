-- name: TrainerFindById :one
SELECT * FROM trainers
WHERE id = ? LIMIT 1;
