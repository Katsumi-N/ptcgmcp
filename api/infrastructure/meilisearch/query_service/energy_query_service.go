package queryservice

import (
	"api/application/search/energy"
	"api/config"
	"context"
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
)

type EnergyMSResponse struct {
	Hits []EnergyResponse `json:"hits"`
}

type EnergyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}

type energyQueryService struct{}

func NewEnergyQueryService() *energyQueryService {
	return &energyQueryService{}
}

func (s *energyQueryService) SearchEnergyList(ctx context.Context, q string) ([]*energy.SearchEnergyList, error) {
	cnf := config.GetConfig()
	msurl := fmt.Sprintf("%s://%s:%s", cnf.MeiliConfig.Protocol, cnf.MeiliConfig.Host, cnf.MeiliConfig.Port)
	client := meilisearch.New(msurl, meilisearch.WithAPIKey(cnf.MeiliConfig.ApiKey))
	index := client.Index("energies")

	searchRes, err := index.Search(q, &meilisearch.SearchRequest{
		Limit: 10,
		Sort:  []string{"id:desc"},
	})
	if err != nil {
		return nil, err
	}

	energyList := lo.Map(searchRes.Hits, func(hit interface{}, _ int) *energy.SearchEnergyList {
		var energyRes EnergyResponse
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(hitBytes, &energyRes); err != nil {
			return nil
		}
		return &energy.SearchEnergyList{
			ID:          energyRes.ID,
			Name:        energyRes.Name,
			ImageURL:    energyRes.ImageURL,
			Description: energyRes.Description,
		}
	})

	energyList = lo.Filter(energyList, func(energy *energy.SearchEnergyList, _ int) bool {
		return energy != nil
	})

	return energyList, nil
}
