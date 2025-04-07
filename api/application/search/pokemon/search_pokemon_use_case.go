package pokemon

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

type SearchPokemonUseCase struct {
	pokemonQueryService PokemonQueryService
}

type SearchPokemonUseCaseDto struct {
	ID         string
	Name       string
	EnergyType string
	Hp         int
	ImageURL   string
	Attacks    []*AttackDto
}

type AttackDto struct {
	Name           string `json:"name"`
	RequiredEnergy string `json:"required_energy"`
	Damage         string `json:"damage"`
	Description    string `json:"description"`
}

func (uc *SearchPokemonUseCase) SearchPokemonList(ctx context.Context, q string) ([]*SearchPokemonUseCaseDto, error) {
	searchPokemonLists, err := uc.pokemonQueryService.SearchPokemonList(ctx, q)
	if err != nil {
		return nil, err
	}

	dtoList := lo.Map(searchPokemonLists, func(f *SearchPokemonList, _ int) *SearchPokemonUseCaseDto {
		attacks := lo.Map(f.Attacks, func(attack PokemonAttackResult, _ int) *AttackDto {
			return &AttackDto{
				Name:           attack.Name,
				RequiredEnergy: attack.RequiredEnergy,
				Damage:         attack.Damage,
				Description:    attack.Description,
			}
		})

		return &SearchPokemonUseCaseDto{
			ID:         fmt.Sprintf("%v", f.ID),
			Name:       f.Name,
			EnergyType: f.EnergyType,
			Hp:         f.Hp,
			ImageURL:   f.ImageURL,
			Attacks:    attacks,
		}
	})

	return dtoList, nil
}
