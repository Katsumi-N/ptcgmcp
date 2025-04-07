package energy

import (
	"api/domain"
	"errors"

	"github.com/samber/lo"
)

// 2025/3/27現在エーススペックのエネルギーはこれらのみ
var aceSpecEnergies = []string{
	"レガシーエネルギー",
	"リッチエネルギー",
	"ネオアッパーエネルギー",
}

type Energy struct {
	id         int
	name       string
	imageUrl   string
	regulation string
	expansion  string
	acespec    bool
}

func NewEnergy(
	id int,
	name string,
	imageUrl string,
	regulation string,
	expansion string,
) (*Energy, error) {
	acespec := isAceSpec(name)
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	if imageUrl == "" {
		return nil, errors.New("image url is required")
	}

	return &Energy{
		id:         id,
		name:       name,
		imageUrl:   imageUrl,
		regulation: regulation,
		expansion:  expansion,
		acespec:    acespec,
	}, nil
}

func isAceSpec(name string) bool {
	return lo.Contains(aceSpecEnergies, name)
}

func (e *Energy) GetId() int {
	return e.id
}

func (e *Energy) GetName() string {
	return e.name
}

func (e *Energy) GetImageUrl() string {
	return e.imageUrl
}

func (e *Energy) GetRegulation() string {
	return e.regulation
}

func (e *Energy) GetExpansion() string {
	return e.expansion
}

func (e *Energy) GetCardType() int {
	return int(domain.Energy)
}

func (e *Energy) IsAceSpec() bool {
	return e.acespec
}
