package deck

import (
	"api/domain"
	domainDeck "api/domain/deck"
	"context"
	"errors"
)

type ICreateDeckUseCase interface {
	Execute(ctx context.Context, request *CreateDeckRequestDto) (*DeckDto, error)
}

type CreateDeckUseCase struct {
	deckRepository domainDeck.DeckRepository
	cardRepository domainDeck.CardRepository
}

func NewCreateDeckUseCase(deckRepository domainDeck.DeckRepository, cardRepository domainDeck.CardRepository) *CreateDeckUseCase {
	return &CreateDeckUseCase{
		deckRepository: deckRepository,
		cardRepository: cardRepository,
	}
}

type CreateDeckRequestDto struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	MainCardID  *CardIDDto           `json:"main_card,omitempty"`
	SubCardID   *CardIDDto           `json:"sub_card,omitempty"`
	Cards       []DeckCardRequestDto `json:"cards"`
}

type CardIDDto struct {
	Id       int    `json:"id"`
	Category string `json:"category"`
}

type DeckCardRequestDto struct {
	Id       int    `json:"id"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
}

func (u *CreateDeckUseCase) Execute(ctx context.Context, request *CreateDeckRequestDto) (*DeckDto, error) {
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

	// デッキの作成
	deck, errs := domainDeck.NewDeck(0, request.Name, request.Description, mainCard, subCard, deckCards)
	if errs != nil {
		errorMsg := ""
		for _, e := range errs {
			errorMsg += e.Error() + "; "
		}
		return nil, errors.New(errorMsg)
	}

	// リポジトリに保存
	createdDeck, err := u.deckRepository.Create(ctx, deck)
	if err != nil {
		return nil, err
	}

	// レスポンス用のDTOを作成
	// メインカードとサブカードのDTO作成
	var mainCardDto *CardDto
	var subCardDto *CardDto

	if createdDeck.GetMainCard() != nil {
		mainCardDto = &CardDto{
			ID:       createdDeck.GetMainCard().GetId(),
			Name:     createdDeck.GetMainCard().GetName(),
			Category: getCardCategory(createdDeck.GetMainCard().GetCardType()),
			ImageURL: createdDeck.GetMainCard().GetImageUrl(),
		}
	}

	if createdDeck.GetSubCard() != nil {
		subCardDto = &CardDto{
			ID:       createdDeck.GetSubCard().GetId(),
			Name:     createdDeck.GetSubCard().GetName(),
			Category: getCardCategory(createdDeck.GetSubCard().GetCardType()),
			ImageURL: createdDeck.GetSubCard().GetImageUrl(),
		}
	}

	// デッキカードの変換
	var deckCardDtos []DeckCardWithQtyDto
	for _, c := range createdDeck.GetCards() {
		deckCardDtos = append(deckCardDtos, DeckCardWithQtyDto{
			ID:       c.GetCard().GetId(),
			Name:     c.GetCard().GetName(),
			Category: getCardCategory(c.GetCard().GetCardType()),
			ImageURL: c.GetCard().GetImageUrl(),
			Quantity: c.GetQuantity(),
		})
	}

	return &DeckDto{
		ID:          createdDeck.GetId(),
		Name:        createdDeck.GetName(),
		Description: createdDeck.GetDescription(),
		MainCard:    mainCardDto,
		SubCard:     subCardDto,
		Cards:       deckCardDtos,
	}, nil
}

// カードタイプをフロントエンド用の文字列に変換
func getCardCategory(cardType interface{}) string {
	switch cardType {
	case 1:
		return "pokemon"
	case 2:
		return "trainer"
	case 3:
		return "energy"
	default:
		return "unknown"
	}
}
