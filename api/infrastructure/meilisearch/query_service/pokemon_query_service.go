package queryservice

import (
	"api/application/search/pokemon"
	"api/config"
	"api/infrastructure/meilisearch/query_service/util"
	"context"
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
)

type PokemonMSResponse struct {
	Hits []PokemonResponse `json:"hits"`
}

type PokemonResponse struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	EnergyType         string          `json:"energy_type"`
	ImageURL           string          `json:"image_url"`
	HP                 int             `json:"hp"`
	Ability            *string         `json:"ability"`
	AbilityDescription *string         `json:"ability_description"`
	Regulation         string          `json:"regulation"`
	Expansion          string          `json:"expansion"`
	Attacks            []PokemonAttack `json:"attacks"`
}

type PokemonAttack struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	RequiredEnergy string  `json:"required_energy"`
	Damage         *string `json:"damage"`
	Description    *string `json:"description"`
}

type pokemonQueryService struct{}

func NewPokemonQueryService() *pokemonQueryService {
	return &pokemonQueryService{}
}

func (s *pokemonQueryService) SearchPokemonList(ctx context.Context, q string) ([]*pokemon.SearchPokemonList, error) {
	cnf := config.GetConfig()
	msurl := fmt.Sprintf("%s://%s:%s", cnf.MeiliConfig.Protocol, cnf.MeiliConfig.Host, cnf.MeiliConfig.Port)
	client := meilisearch.New(msurl, meilisearch.WithAPIKey(cnf.MeiliConfig.ApiKey))
	index := client.Index("pokemons")

	q = util.HiraganaToKatakana(q)
	searchRes, err := index.Search(q, &meilisearch.SearchRequest{
		Limit: 10,
		Sort:  []string{"id:desc"},
	})
	if err != nil {
		return nil, err
	}

	pokemonList := lo.Map(searchRes.Hits, func(hit interface{}, _ int) *pokemon.SearchPokemonList {
		var pokemon pokemon.SearchPokemonList
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(hitBytes, &pokemon); err != nil {
			return nil
		}
		return &pokemon
	})

	return pokemonList, nil
}
