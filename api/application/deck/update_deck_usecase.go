package deck

import (
	"api/domain"
	domainDeck "api/domain/deck"
	"context"
	"errors"
	"fmt"
)

type IUpdateDeckUseCase interface {
	Execute(ctx context.Context, id int, request *UpdateDeckRequestDto) (*DeckDto, error)
}

type UpdateDeckUseCase struct {
	deckRepository domainDeck.DeckRepository
	cardRepository domainDeck.CardRepository
}

func NewUpdateDeckUseCase(deckRepository domainDeck.DeckRepository, cardRepository domainDeck.CardRepository) *UpdateDeckUseCase {
	return &UpdateDeckUseCase{
		deckRepository: deckRepository,
		cardRepository: cardRepository,
	}
}

type UpdateDeckRequestDto struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	MainCardID  *CardIDDto           `json:"main_card,omitempty"`
	SubCardID   *CardIDDto           `json:"sub_card,omitempty"`
	Cards       []DeckCardRequestDto `json:"cards"`
}

func (u *UpdateDeckUseCase) Execute(ctx context.Context, id int, request *UpdateDeckRequestDto) (*DeckDto, error) {
	// 既存デッキを取得
	_, err := u.deckRepository.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("デッキが見つかりません: %w", err)
	}

	// カード情報の取得
	var mainCard domain.Card
	var subCard domain.Card
	var deckCards []domainDeck.DeckCard

	// メインカードの取得（存在する場合）
	if request.MainCardID != nil {
		cardType, exists := domain.StringToCardType[request.MainCardID.Category]
		if !exists {
			return nil, errors.New("メインカードのカテゴリが不正です")
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
			return nil, errors.New("サブカードのカテゴリが不正です")
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
			return nil, errors.New("カードのカテゴリが不正です")
		}
		card, err := u.cardRepository.FindCardById(ctx, cardRequest.Id, cardType)
		if err != nil {
			return nil, err
		}
		deckCard := domainDeck.NewDeckCard(card, cardRequest.Quantity)
		deckCards = append(deckCards, *deckCard)
	}

	// デッキの作成（既存IDを保持）
	deck, errs := domainDeck.NewDeck(id, request.Name, request.Description, mainCard, subCard, deckCards)
	if errs != nil {
		errorMsg := ""
		for _, e := range errs {
			errorMsg += e.Error() + "; "
		}
		return nil, errors.New(errorMsg)
	}

	// リポジトリで更新
	if err := u.deckRepository.Update(ctx, deck); err != nil {
		return nil, err
	}

	// 更新後のデッキを再取得
	updatedDeck, err := u.deckRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// レスポンス用のDTOを作成
	// メインカードとサブカードのDTO作成
	var mainCardDto *CardDto
	var subCardDto *CardDto

	if updatedDeck.GetMainCard() != nil {
		mainCardDto = &CardDto{
			ID:       updatedDeck.GetMainCard().GetId(),
			Name:     updatedDeck.GetMainCard().GetName(),
			Category: getCardCategory(updatedDeck.GetMainCard().GetCardType()),
			ImageURL: updatedDeck.GetMainCard().GetImageUrl(),
		}
	}

	if updatedDeck.GetSubCard() != nil {
		subCardDto = &CardDto{
			ID:       updatedDeck.GetSubCard().GetId(),
			Name:     updatedDeck.GetSubCard().GetName(),
			Category: getCardCategory(updatedDeck.GetSubCard().GetCardType()),
			ImageURL: updatedDeck.GetSubCard().GetImageUrl(),
		}
	}

	// デッキカードの変換
	var deckCardDtos []DeckCardWithQtyDto
	for _, c := range updatedDeck.GetCards() {
		deckCardDtos = append(deckCardDtos, DeckCardWithQtyDto{
			ID:       c.GetCard().GetId(),
			Name:     c.GetCard().GetName(),
			Category: getCardCategory(c.GetCard().GetCardType()),
			ImageURL: c.GetCard().GetImageUrl(),
			Quantity: c.GetQuantity(),
		})
	}

	return &DeckDto{
		ID:          updatedDeck.GetId(),
		Name:        updatedDeck.GetName(),
		Description: updatedDeck.GetDescription(),
		MainCard:    mainCardDto,
		SubCard:     subCardDto,
		Cards:       deckCardDtos,
	}, nil
}
