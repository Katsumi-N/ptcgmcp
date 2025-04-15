package deck

import "context"

type SearchDeckListDto struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	MainCard    SearchDeckCardDto   `json:"main_card"`
	SubCard     SearchDeckCardDto   `json:"sub_card"`
	Cards       []SearchDeckCardDto `json:"cards"`
}

type SearchDeckCardDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	ImageURL string `json:"image_url"`
}

type DeckQueryService interface {
	SearchDeckList(ctx context.Context, q string) ([]*SearchDeckListDto, error)
}
