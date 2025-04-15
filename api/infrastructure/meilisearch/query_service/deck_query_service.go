package queryservice

import (
	"api/application/search/deck"
	"api/config"
	"context"
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
)

type deckQueryService struct{}

func NewDeckQueryService() *deckQueryService {
	return &deckQueryService{}
}

type DeckMSResponse struct {
	Hits []DeckResponse `json:"hits"`
}

type DeckResponse struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	MainCard    DeckCardResponse   `json:"main_card"`
	SubCard     DeckCardResponse   `json:"sub_card"`
	Cards       []DeckCardResponse `json:"cards"`
}

type DeckCardResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	ImageURL string `json:"image_url"`
}

func (d *deckQueryService) SearchDeckList(ctx context.Context, q string) ([]*deck.SearchDeckListDto, error) {
	cnf := config.GetConfig()
	msurl := cnf.MeiliConfig.Protocol + "://" + cnf.MeiliConfig.Host + ":" + cnf.MeiliConfig.Port
	client := meilisearch.New(msurl, meilisearch.WithAPIKey(cnf.MeiliConfig.ApiKey))
	index := client.Index("decks")
	searchRes, err := index.Search(q, &meilisearch.SearchRequest{
		Limit: 10,
		Sort:  []string{"id:desc"},
	})
	if err != nil {
		return nil, err
	}

	deckList := lo.Map(searchRes.Hits, func(hit interface{}, _ int) *deck.SearchDeckListDto {
		var deckRes DeckResponse
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(hitBytes, &deckRes); err != nil {
			return nil
		}
		return &deck.SearchDeckListDto{
			Id:          deckRes.ID,
			Name:        deckRes.Name,
			Description: deckRes.Description,
			MainCard: deck.SearchDeckCardDto{
				Id:       deckRes.MainCard.ID,
				Name:     deckRes.MainCard.Name,
				Category: deckRes.MainCard.Category,
				Quantity: deckRes.MainCard.Quantity,
				ImageURL: deckRes.MainCard.ImageURL,
			},
			SubCard: deck.SearchDeckCardDto{
				Id:       deckRes.SubCard.ID,
				Name:     deckRes.SubCard.Name,
				Category: deckRes.SubCard.Category,
				Quantity: deckRes.SubCard.Quantity,
				ImageURL: deckRes.SubCard.ImageURL,
			},
			Cards: lo.Map(deckRes.Cards, func(card DeckCardResponse, _ int) deck.SearchDeckCardDto {
				return deck.SearchDeckCardDto{
					Id:       card.ID,
					Name:     card.Name,
					Category: card.Category,
					Quantity: card.Quantity,
					ImageURL: card.ImageURL,
				}
			}),
		}
	})

	return deckList, nil
}
