package deck

import (
	"api/domain"
	domainDeck "api/domain/deck"
	"context"
	"errors"
)

type IValidateDeckUseCase interface {
	Execute(ctx context.Context, request *ValidateDeckRequestDto) (*ValidateDeckResponseDto, error)
}

type ValidateDeckUseCase struct {
	cardRepository domainDeck.CardRepository
}

func NewValidateDeckUseCase(cardRepository domainDeck.CardRepository) *ValidateDeckUseCase {
	return &ValidateDeckUseCase{
		cardRepository: cardRepository,
	}
}

type ValidateDeckRequestDto struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	MainCardID  *CardIDDto           `json:"main_card,omitempty"`
	SubCardID   *CardIDDto           `json:"sub_card,omitempty"`
	Cards       []DeckCardRequestDto `json:"cards"`
}

type ValidateDeckResponseDto struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

func (u *ValidateDeckUseCase) Execute(ctx context.Context, request *ValidateDeckRequestDto) (*ValidateDeckResponseDto, error) {
	// カード情報の取得
	var mainCard domain.Card
	var subCard domain.Card
	var deckCards []domainDeck.DeckCard

	// メインカードの取得（存在する場合）
	if request.MainCardID != nil {
		cardType, exists := domain.StringToCardType[request.MainCardID.Category]
		if !exists {
			return nil, errors.New("invalid main card category")
		}
		card, err := u.cardRepository.FindCardById(ctx, request.MainCardID.Id, cardType)
		if err != nil {
			return nil, err
		}
		mainCard = card
	}

	// サブカードの取得（存在する場合）
	if request.SubCardID != nil {
		cardType, exists := domain.StringToCardType[request.SubCardID.Category]
		if !exists {
			return nil, errors.New("invalid sub card category")
		}
		card, err := u.cardRepository.FindCardById(ctx, request.SubCardID.Id, cardType)
		if err != nil {
			return nil, err
		}
		subCard = card
	}

	// デッキカードの取得
	for _, cardRequest := range request.Cards {
		cardType, exists := domain.StringToCardType[cardRequest.Category]
		if !exists {
			return nil, errors.New("invalid card category")
		}
		card, err := u.cardRepository.FindCardById(ctx, cardRequest.Id, cardType)
		if err != nil {
			return nil, err
		}

		deckCard := domainDeck.NewDeckCard(card, cardRequest.Quantity)
		deckCards = append(deckCards, *deckCard)
	}

	_, validationErrors := domainDeck.NewDeck(0, request.Name, request.Description, mainCard, subCard, deckCards)

	// エラーメッセージをレスポンス用に変換
	var errorMessages []string
	for _, e := range validationErrors {
		errorMessages = append(errorMessages, e.Error())
	}

	return &ValidateDeckResponseDto{
		IsValid: len(errorMessages) == 0,
		Errors:  errorMessages,
	}, nil
}
