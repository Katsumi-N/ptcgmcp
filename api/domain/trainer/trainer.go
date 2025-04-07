package trainer

import (
	"api/domain"
	"errors"

	"github.com/samber/lo"
)

type Trainer struct {
	id          int
	name        string
	trainerType string
	description string
	imageUrl    string
	regulation  string
	expansion   string
	acespec     bool
}

const (
	Supporter           = "サポート"
	Stadium             = "スタジアム"
	Item                = "グッズ"
	PokemonsItem        = "ポケモンのどうぐ"
	AceSpecItem         = "グッズ特別なルール"
	AceSpecPokemonsItem = "ポケモンのどうぐ特別なルール"
	AceSpecStadium      = "スタジアム特別なルール"
)

var validTrainerTypes = []string{Supporter, Stadium, Item, PokemonsItem}
var aceSpecTrainerTypes = []string{AceSpecItem, AceSpecPokemonsItem, AceSpecStadium}

func NewTrainer(id int, name string, trainerType string, description string, imageUrl string, regulation string, expansion string) (*Trainer, error) {
	// エーススペックを確認
	// TODO: テーブルにacespec カラム追加
	acespec := isAceSpec(trainerType)

	if !isValidTrainerType(trainerType) {
		return nil, errors.New("Trainer type must be supporter, stadium or item")
	}

	return &Trainer{
		id:          id,
		name:        name,
		trainerType: trainerType,
		description: description,
		imageUrl:    imageUrl,
		regulation:  regulation,
		expansion:   expansion,
		acespec:     acespec,
	}, nil
}

func isValidTrainerType(trainerType string) bool {
	return lo.Contains(validTrainerTypes, trainerType)
}

func isAceSpec(trainerType string) bool {
	return lo.Contains(aceSpecTrainerTypes, trainerType)
}

func (t *Trainer) GetId() int {
	return t.id
}

func (t *Trainer) GetName() string {
	return t.name
}

func (t *Trainer) GetCardType() int {
	return int(domain.Trainer)
}

func (t *Trainer) GetImageUrl() string {
	return t.imageUrl
}

func (t *Trainer) IsAceSpec() bool {
	return t.acespec
}
