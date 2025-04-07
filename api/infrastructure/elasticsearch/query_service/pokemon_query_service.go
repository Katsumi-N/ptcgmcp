package query_service

import (
	"api/application/search/pokemon"
	"api/config"
	"api/infrastructure/elasticsearch/query_service/util"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

// Elasticsearchから返されるPokemonのJSON構造体
type PokemonESResponse struct {
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

func NewPokemonQueryService() pokemon.PokemonQueryService {
	return &pokemonQueryService{}
}

func (s *pokemonQueryService) SearchPokemonList(ctx context.Context, q string) ([]*pokemon.SearchPokemonList, error) {
	cnf := config.GetConfig()
	esUrl := fmt.Sprintf("%s://%s:%s", cnf.ESConfig.EsProtocol, cnf.ESConfig.EsHost, cnf.ESConfig.EsPort)
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			esUrl,
		},
	})
	if err != nil {
		log.Println("error creating elasticsearch client: ", err)
		return nil, err
	}

	q = util.HiraganaToKatakana(q)

	// クエリの構築
	req := &search.Request{
		Size: util.IntPtr(100),
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Should: []types.Query{
					{
						Match: map[string]types.MatchQuery{
							"name": {Query: q, Boost: util.Float32Ptr(3.0)},
						},
					},
				},
			},
		},
		Sort: []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"id": {Order: &sortorder.Desc},
				}},
		},
	}

	res, err := es.Search().Index("pokemons").Request(req).Do(ctx)
	if err != nil {
		log.Println("error searching elasticsearch: ", err)
		return nil, err
	}

	var searchPokemonList []*pokemon.SearchPokemonList
	for _, hit := range res.Hits.Hits {
		var esResponse PokemonESResponse
		if err := json.Unmarshal(hit.Source_, &esResponse); err != nil {
			log.Println("error unmarshalling hit source: ", err)
			return nil, err
		}

		// PokemonESResponseからSearchPokemonListへの変換
		attackResults := make([]pokemon.PokemonAttackResult, 0, len(esResponse.Attacks))
		for _, attack := range esResponse.Attacks {
			damage := ""
			description := ""
			if attack.Damage != nil {
				damage = *attack.Damage
			}
			if attack.Description != nil {
				description = *attack.Description
			}
			attackResults = append(attackResults, pokemon.PokemonAttackResult{
				Name:           attack.Name,
				RequiredEnergy: attack.RequiredEnergy,
				Damage:         damage,
				Description:    description,
			})
		}

		searchPokemon := &pokemon.SearchPokemonList{
			ID:         esResponse.ID,
			Name:       esResponse.Name,
			EnergyType: esResponse.EnergyType,
			Hp:         esResponse.HP,
			ImageURL:   esResponse.ImageURL,
			Attacks:    attackResults,
		}
		searchPokemonList = append(searchPokemonList, searchPokemon)
	}

	return searchPokemonList, nil
}
