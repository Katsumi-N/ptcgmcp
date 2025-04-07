package query_service

import (
	"api/application/detail"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetailQueryService_FindPokemonDetail(t *testing.T) {
	setupFixtures(t)

	repo := NewDetailQueryService()
	attacks := []detail.PokemonAttack{
		{
			Name:           "トパーズボルト",
			RequiredEnergy: "草雷鋼",
			Damage:         "300",
			Description:    "このポケモンについているエネルギーを3個選び、トラッシュする。",
		},
	}
	pika := &detail.Pokemon{
		Id:                 1,
		Name:               "ピカチュウex",
		EnergyType:         "雷",
		Hp:                 200,
		Ability:            "",
		AbilityDescription: "",
		ImageUrl:           "pika.png",
		Regulation:         "",
		Expansion:          "",
		Attacks:            attacks,
	}
	tests := map[string]struct {
		pokemonId   int
		expected    *detail.Pokemon
		expectError bool
	}{
		"valid_pokemon": {
			pokemonId:   1,
			expected:    pika,
			expectError: false,
		},
		"invalid_pokemon": {
			pokemonId:   999,
			expected:    nil,
			expectError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			result, err := repo.FindPokemonDetail(ctx, tt.pokemonId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
