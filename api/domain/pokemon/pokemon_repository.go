package pokemon

import (
	"context"
)

type PokemonRepository interface {
	FindById(ctx context.Context, pokemonId int) (*Pokemon, error)
}
