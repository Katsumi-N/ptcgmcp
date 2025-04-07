package detail

import "context"

type Pokemon struct {
	Id                 int
	Name               string
	EnergyType         string
	ImageUrl           string
	Hp                 int
	Ability            string
	AbilityDescription string
	Regulation         string
	Expansion          string
	Attacks            []PokemonAttack
}

type PokemonAttack struct {
	Name           string
	RequiredEnergy string
	Damage         string
	Description    string
}

type Trainer struct {
	Id          int
	Name        string
	TrainerType string
	Description string
	ImageUrl    string
	Regulation  string
	Expansion   string
}

type Energy struct {
	Id          int
	Name        string
	ImageUrl    string
	Description string
	Regulation  string
	Expansion   string
}

type DetailQueryService interface {
	FindPokemonDetail(ctx context.Context, pokemonId int) (*Pokemon, error)
	FindTrainerDetail(ctx context.Context, trainerId int) (*Trainer, error)
	FindEnergyDetail(ctx context.Context, energyId int) (*Energy, error)
}
