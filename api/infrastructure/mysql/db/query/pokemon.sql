-- name: PokemonFindById :one
SELECT * FROM pokemons
WHERE id = ? LIMIT 1;

-- name: PokemonAttackFindByPokemonId :many
SELECT * FROM pokemon_attacks
WHERE pokemon_id = ?;
