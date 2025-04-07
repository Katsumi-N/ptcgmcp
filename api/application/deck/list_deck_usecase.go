package deck

import (
	"api/domain/deck"
	"context"
	"errors"
)

type IListDeckUseCase interface {
	GetAllDecks(ctx context.Context) ([]*DeckDto, error)
	GetDeckById(ctx context.Context, deckId int) (*DeckDto, error)
}

type ListDeckUseCase struct {
	deckRepository deck.DeckRepository
}

func NewListDeckUseCase(deckRepository deck.DeckRepository) *ListDeckUseCase {
	return &ListDeckUseCase{
		deckRepository: deckRepository,
	}
}

func (u *ListDeckUseCase) GetAllDecks(ctx context.Context) ([]*DeckDto, error) {
	decks, err := u.deckRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var deckDtos []*DeckDto
	for _, d := range decks {
		var mainCardDto *CardDto
		var subCardDto *CardDto

		// メインカードがある場合は変換
		if d.GetMainCard() != nil {
			mainCardDto = &CardDto{
				ID:       d.GetMainCard().GetId(),
				Name:     d.GetMainCard().GetName(),
				Category: getCardCategory(d.GetMainCard().GetCardType()),
				ImageURL: d.GetMainCard().GetImageUrl(),
			}
		}

		// サブカードがある場合は変換
		if d.GetSubCard() != nil {
			subCardDto = &CardDto{
				ID:       d.GetSubCard().GetId(),
				Name:     d.GetSubCard().GetName(),
				Category: getCardCategory(d.GetSubCard().GetCardType()),
				ImageURL: d.GetSubCard().GetImageUrl(),
			}
		}

		// デッキカードの変換
		var deckCardDtos []DeckCardWithQtyDto
		for _, c := range d.GetCards() {
			deckCardDtos = append(deckCardDtos, DeckCardWithQtyDto{
				ID:       c.GetCard().GetId(),
				Name:     c.GetCard().GetName(),
				Category: getCardCategory(c.GetCard().GetCardType()),
				ImageURL: c.GetCard().GetImageUrl(),
				Quantity: c.GetQuantity(),
			})
		}

		deckDtos = append(deckDtos, &DeckDto{
			ID:          d.GetId(),
			Name:        d.GetName(),
			Description: d.GetDescription(),
			MainCard:    mainCardDto,
			SubCard:     subCardDto,
			Cards:       deckCardDtos,
		})
	}

	return deckDtos, nil
}

func (u *ListDeckUseCase) GetDeckById(ctx context.Context, deckId int) (*DeckDto, error) {
	d, err := u.deckRepository.FindById(ctx, deckId)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, errors.New("デッキが見つかりません")
	}

	var mainCardDto *CardDto
	var subCardDto *CardDto
	// メインカードがある場合は変換
	if d.GetMainCard() != nil {
		mainCardDto = &CardDto{
			ID:       d.GetMainCard().GetId(),
			Name:     d.GetMainCard().GetName(),
			Category: getCardCategory(d.GetMainCard().GetCardType()),
			ImageURL: d.GetMainCard().GetImageUrl(),
		}
	}
	// サブカードがある場合は変換
	if d.GetSubCard() != nil {
		subCardDto = &CardDto{
			ID:       d.GetSubCard().GetId(),
			Name:     d.GetSubCard().GetName(),
			Category: getCardCategory(d.GetSubCard().GetCardType()),
			ImageURL: d.GetSubCard().GetImageUrl(),
		}
	}

	// デッキカードの変換
	var deckCardDtos []DeckCardWithQtyDto
	for _, c := range d.GetCards() {
		deckCardDtos = append(deckCardDtos, DeckCardWithQtyDto{
			ID:       c.GetCard().GetId(),
			Name:     c.GetCard().GetName(),
			Category: getCardCategory(c.GetCard().GetCardType()),
			ImageURL: c.GetCard().GetImageUrl(),
			Quantity: c.GetQuantity(),
		})
	}

	return &DeckDto{
		ID:          d.GetId(),
		Name:        d.GetName(),
		Description: d.GetDescription(),
		MainCard:    mainCardDto,
		SubCard:     subCardDto,
		Cards:       deckCardDtos,
	}, nil
}
