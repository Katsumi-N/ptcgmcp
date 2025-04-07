package pokemon

import (
	"testing"
)

func TestNewPokemon(t *testing.T) {
	tests := map[string]struct {
		id                 int
		name               string
		energyType         string
		hp                 int
		ability            string
		abilityDescription string
		description        string
		imageUrl           string
		regulation         string
		expansion          string
		acespec            bool
		attacks            []PokemonAttack
		expectError        bool
	}{
		"valid": {
			id:                 1,
			name:               "ピカチュウ",
			energyType:         Electric,
			hp:                 120,
			ability:            "Static",
			abilityDescription: "May cause paralysis if touched.",
			description:        "ぴかぴか",
			imageUrl:           "https://example.com/pikachu.png",
			regulation:         "Regulation A",
			expansion:          "Expansion 1",
			attacks:            []PokemonAttack{NewPokemonAttack("Thunder Shock", Electric, "30", "A jolt of electricity")},
			expectError:        false,
		},
		"invalid energy type": {
			id:                 2,
			name:               "ラルトス",
			energyType:         "Fairy",
			hp:                 120,
			ability:            "Synchronize",
			abilityDescription: "Passes a burn, poison, or paralysis to the foe.",
			description:        "フェアリーがなくなった",
			imageUrl:           "https://example.com/ralts.png",
			regulation:         "Regulation B",
			expansion:          "Expansion 2",
			attacks:            []PokemonAttack{NewPokemonAttack("Confusion", Psychic, "20", "Confuses the opponent")},
			expectError:        true,
		},
		"invalid hp": {
			id:                 3,
			name:               "ホゲータ",
			energyType:         Fire,
			hp:                 -10,
			ability:            "Blaze",
			abilityDescription: "Powers up Fire-type moves when the Pokémon's HP is low.",
			description:        "ほげほげ",
			imageUrl:           "https://example.com/hoge.png",
			regulation:         "Regulation C",
			expansion:          "Expansion 3",
			attacks:            []PokemonAttack{NewPokemonAttack("Ember", Fire, "40", "A small flame")},
			expectError:        true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pokemon, err := NewPokemon(tt.id, tt.name, tt.energyType, tt.hp, tt.ability, tt.abilityDescription, tt.imageUrl, tt.regulation, tt.expansion, tt.attacks)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for input %v, got nil", tt)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %v: %v", tt, err)
				}
				if pokemon.id != tt.id || pokemon.name != tt.name || pokemon.energyType != tt.energyType || pokemon.hp != tt.hp || pokemon.ability != tt.ability || pokemon.abilityDescription != tt.abilityDescription || pokemon.imageUrl != tt.imageUrl || pokemon.regulation != tt.regulation || pokemon.expansion != tt.expansion || len(pokemon.attacks) != len(tt.attacks) {
					t.Errorf("expected %v, got %v", tt, pokemon)
				}
			}
		})
	}
}
