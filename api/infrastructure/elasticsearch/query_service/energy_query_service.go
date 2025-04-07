package query_service

import (
	"api/application/search/energy"
	"api/config"
	"api/infrastructure/elasticsearch/query_service/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

// Elasticsearchから返されるEnergyのJSON構造体
type EnergyESResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Regulation  string    `json:"regulation"`
	Expansion   string    `json:"expansion"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type energyQueryService struct{}

func NewEnergyQueryService() energy.EnergyQueryService {
	return &energyQueryService{}
}

func (s *energyQueryService) SearchEnergyList(ctx context.Context, q string) ([]*energy.SearchEnergyList, error) {
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

	// 元のクエリqと、ひらがなをカタカナに変換したクエリを両方使用
	originalQ := q
	katakanaQ := util.HiraganaToKatakana(q)

	// クエリの構築
	req := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Should: []types.Query{
					{
						Match: map[string]types.MatchQuery{
							"name": {Query: originalQ, Boost: util.Float32Ptr(2.0)},
						},
					},
					{
						Match: map[string]types.MatchQuery{
							"name": {Query: katakanaQ, Boost: util.Float32Ptr(2.0)},
						},
					},
					{
						Match: map[string]types.MatchQuery{
							"description": {Query: originalQ},
						},
					},
				},
				MinimumShouldMatch: util.StringPtr("1"),
			},
		},
		Sort: []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"id": {Order: &sortorder.Desc},
				}},
		},
	}
	res, err := es.Search().Index("energies").Request(req).Do(ctx)
	if err != nil {
		log.Println("error searching elasticsearch: ", err)
		return nil, err
	}

	var searchEnergyList []*energy.SearchEnergyList
	for _, hit := range res.Hits.Hits {
		var esResponse EnergyESResponse
		if err := json.Unmarshal(hit.Source_, &esResponse); err != nil {
			log.Println("error unmarshalling hit source: ", err)
			return nil, err
		}

		// EnergyESResponseからSearchEnergyListへの変換
		searchEnergy := &energy.SearchEnergyList{
			ID:          esResponse.ID,
			Name:        esResponse.Name,
			ImageURL:    esResponse.ImageURL,
			Description: esResponse.Description,
		}
		searchEnergyList = append(searchEnergyList, searchEnergy)
	}

	return searchEnergyList, nil
}
