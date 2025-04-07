package pokemon

import (
	"api/domain"
	"errors"

	"github.com/samber/lo"
)

type Pokemon struct {
	id                 int
	name               string
	energyType         string
	imageUrl           string
	hp                 int
	ability            string
	abilityDescription string
	regulation         string
	expansion          string
	attacks            []PokemonAttack
	acespec            bool
}

type PokemonAttack struct {
	name           string
	requiredEnergy string
	damage         string
	description    string
}

const (
	Fire      = "炎"
	Water     = "水"
	Electric  = "雷"
	Fight     = "闘"
	Psychic   = "超"
	Grass     = "草"
	Steel     = "鋼"
	Dark      = "悪"
	Colorless = "無色"
	Dragon    = "ドラゴン"
)

var validEnergyTypes = []string{Fire, Water, Electric, Fight, Psychic, Grass, Steel, Dark, Colorless, Dragon}

func NewPokemonAttack(name string, requiredEnergy string, damage string, description string) PokemonAttack {
	return PokemonAttack{
		name:           name,
		requiredEnergy: requiredEnergy,
		damage:         damage,
		description:    description,
	}
}

func NewPokemon(id int, name string, energyType string, hp int, ability string, abilityDescription string, imageUrl string, regulation string, expansion string, attacks []PokemonAttack, isAceSpec ...bool) (*Pokemon, error) {
	// 今の所ポケモンにエーススペックはいない
	acespec := false

	if !isValidEnergyType(energyType) {
		return nil, errors.New("energy type must be valid type")
	}

	if hp < 0 {
		return nil, errors.New("HP must be greater than or equal to 0")
	}

	return &Pokemon{
		id:                 id,
		name:               name,
		energyType:         energyType,
		hp:                 hp,
		ability:            ability,
		abilityDescription: abilityDescription,
		imageUrl:           imageUrl,
		regulation:         regulation,
		expansion:          expansion,
		attacks:            attacks,
		acespec:            acespec,
	}, nil
}

func isValidEnergyType(energyType string) bool {
	return lo.Contains(validEnergyTypes, energyType)
}

func (p *Pokemon) GetId() int {
	return p.id
}

func (p *Pokemon) GetName() string {
	return p.name
}

func (p *Pokemon) GetCardType() int {
	return int(domain.Pokemon)
}

func (p *Pokemon) GetImageUrl() string {
	return p.imageUrl
}

func (p *Pokemon) IsAceSpec() bool {
	return p.acespec
}
