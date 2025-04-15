package deck

import (
	"context"

	"github.com/samber/lo"
)

type ISearchDeckUseCase interface {
	SearchDeckList(ctx context.Context, q string) ([]*SearchDeckUseCaseDto, error)
}

type SearchDeckUseCase struct {
	deckQueryService DeckQueryService
}

func NewSearchDeckUseCase(deckQueryService DeckQueryService) *SearchDeckUseCase {
	return &SearchDeckUseCase{
		deckQueryService: deckQueryService,
	}
}

type SearchDeckUseCaseDto struct {
	Id          int                        `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	MainCard    SearchDeckCardUseCaseDto   `json:"main_card"`
	SubCard     SearchDeckCardUseCaseDto   `json:"sub_card"`
	Cards       []SearchDeckCardUseCaseDto `json:"cards"`
}

type SearchDeckCardUseCaseDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	ImageURL string `json:"image_url"`
}

func NewSearchDeckUseCaseDto(id int, name, description string, mainCard, subCard SearchDeckCardUseCaseDto, cards []SearchDeckCardUseCaseDto) *SearchDeckUseCaseDto {
	return &SearchDeckUseCaseDto{
		Id:          id,
		Name:        name,
		Description: description,
		MainCard:    mainCard,
		SubCard:     subCard,
		Cards:       cards,
	}
}
func NewSearchDeckCardUseCaseDto(id int, name, category string, quantity int, imageURL string) *SearchDeckCardUseCaseDto {
	return &SearchDeckCardUseCaseDto{
		Id:       id,
		Name:     name,
		Category: category,
		Quantity: quantity,
		ImageURL: imageURL,
	}
}

func (u *SearchDeckUseCase) SearchDeckList(ctx context.Context, q string) ([]*SearchDeckUseCaseDto, error) {
	decks, err := u.deckQueryService.SearchDeckList(ctx, q)
	if err != nil {
		return nil, err
	}

	deckList := lo.Map(decks, func(f *SearchDeckListDto, _ int) *SearchDeckUseCaseDto {
		card := lo.Map(f.Cards, func(card SearchDeckCardDto, _ int) SearchDeckCardUseCaseDto {
			return *NewSearchDeckCardUseCaseDto(card.Id, card.Name, card.Category, card.Quantity, card.ImageURL)
		})

		return &SearchDeckUseCaseDto{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			MainCard: SearchDeckCardUseCaseDto{
				Id:       f.MainCard.Id,
				Name:     f.MainCard.Name,
				Category: f.MainCard.Category,
				ImageURL: f.MainCard.ImageURL,
			},
			SubCard: SearchDeckCardUseCaseDto{
				Id:       f.SubCard.Id,
				Name:     f.SubCard.Name,
				Category: f.SubCard.Category,
				ImageURL: f.SubCard.ImageURL,
			},
			Cards: card,
		}
	})

	return deckList, nil
}
