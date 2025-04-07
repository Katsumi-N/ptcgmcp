package queryservice

import (
	"api/application/search/trainer"
	"api/config"
	"context"
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
)

type TrainerMSResponse struct {
	Hits []TrainerResponse `json:"hits"`
}

type TrainerResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TrainerType string `json:"trainer_type"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}

type trainerQueryService struct{}

func NewTrainerQueryService() *trainerQueryService {
	return &trainerQueryService{}
}

func (s *trainerQueryService) SearchTrainerList(ctx context.Context, q string) ([]*trainer.SearchTrainerList, error) {
	cnf := config.GetConfig()
	msurl := fmt.Sprintf("%s://%s:%s", cnf.MeiliConfig.Protocol, cnf.MeiliConfig.Host, cnf.MeiliConfig.Port)
	client := meilisearch.New(msurl, meilisearch.WithAPIKey(cnf.MeiliConfig.ApiKey))
	index := client.Index("trainers")

	searchRes, err := index.Search(q, &meilisearch.SearchRequest{
		Limit: 10,
		Sort:  []string{"id:desc"},
	})
	if err != nil {
		return nil, err
	}

	trainerList := lo.Map(searchRes.Hits, func(hit interface{}, _ int) *trainer.SearchTrainerList {
		var trainerRes TrainerResponse
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(hitBytes, &trainerRes); err != nil {
			return nil
		}
		return &trainer.SearchTrainerList{
			ID:          trainerRes.ID,
			Name:        trainerRes.Name,
			TrainerType: trainerRes.TrainerType,
			ImageURL:    trainerRes.ImageURL,
		}
	})

	trainerList = lo.Filter(trainerList, func(trainer *trainer.SearchTrainerList, _ int) bool {
		return trainer != nil
	})

	return trainerList, nil
}
