-- name: CreateDeck :execresult
INSERT INTO decks (
  name,
  description,
  main_card_id,
  main_card_type_id,
  sub_card_id,
  sub_card_type_id
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: CreateDeckCard :execresult
INSERT INTO deck_cards (
  deck_id,
  card_id,
  card_type_id,
  quantity
) VALUES (
  ?, ?, ?, ?
);

-- name: FindALl :many
SELECT * FROM decks
ORDER BY id DESC;

-- name: FindDeckById :one
SELECT * FROM decks
WHERE id = ?
LIMIT 1;

-- name: FindDeckCardsByDeckId :many
SELECT * FROM deck_cards
WHERE deck_id = ?;

-- name: UpdateDeck :exec
UPDATE decks
SET 
  name = ?,
  description = ?,
  main_card_id = ?,
  main_card_type_id = ?,
  sub_card_id = ?,
  sub_card_type_id = ?
WHERE id = ?;

-- name: DeleteDeck :exec
DELETE FROM decks
WHERE id = ?;

-- name: DeleteDeckCardsByDeckId :exec
DELETE FROM deck_cards
WHERE deck_id = ?;