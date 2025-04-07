package deck

import (
	"api/domain"
	domainDeck "api/domain/deck"
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックカード
type mockCard struct {
	id       int
	name     string
	cardType int
	imageUrl string
	aceSpec  bool
}

func (m *mockCard) GetId() int {
	return m.id
}

func (m *mockCard) GetName() string {
	return m.name
}

func (m *mockCard) GetCardType() int {
	return m.cardType
}

func (m *mockCard) GetImageUrl() string {
	return m.imageUrl
}

func (m *mockCard) IsAceSpec() bool {
	return m.aceSpec
}

// モックデッキリポジトリ
type mockDeckRepository struct {
	mock.Mock
}

func (m *mockDeckRepository) Create(ctx context.Context, d *domainDeck.Deck) (*domainDeck.Deck, error) {
	args := m.Called(ctx, d)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainDeck.Deck), args.Error(1)
}

func (m *mockDeckRepository) FindById(ctx context.Context, id int) (*domainDeck.Deck, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainDeck.Deck), args.Error(1)
}

func (m *mockDeckRepository) FindAll(ctx context.Context) ([]*domainDeck.Deck, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainDeck.Deck), args.Error(1)
}

func (m *mockDeckRepository) Update(ctx context.Context, d *domainDeck.Deck) error {
	args := m.Called(ctx, d)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *mockDeckRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// モックカードリポジトリ
type mockCardRepository struct {
	mock.Mock
}

func (m *mockCardRepository) FindCardById(ctx context.Context, id int, cardType domain.CardType) (domain.Card, error) {
	args := m.Called(ctx, id, cardType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Card), args.Error(1)
}

func TestExecute(t *testing.T) {
	// テストケースの定義
	tests := map[string]struct {
		request         *CreateDeckRequestDto
		mockCards       map[string]domain.Card
		returnDeck      *domainDeck.Deck
		expectError     bool
		expectedErrMsg  string
		expectedDeckDto *DeckDto
	}{
		"success": {
			request: &CreateDeckRequestDto{
				Name:        "テストデッキ",
				Description: "テスト用のデッキです",
				MainCardID: &CardIDDto{
					Id:       1,
					Category: "pokemon",
				},
				SubCardID: &CardIDDto{
					Id:       2,
					Category: "trainer",
				},
				Cards: []DeckCardRequestDto{
					{
						Id:       1,
						Category: "pokemon",
						Quantity: 4,
					},
					{
						Id:       2,
						Category: "trainer",
						Quantity: 4,
					},
					{
						Id:       3,
						Category: "energy",
						Quantity: 52,
					},
				},
			},
			mockCards: map[string]domain.Card{
				"pokemon": &mockCard{id: 1, name: "ピカチュウ", cardType: 1, imageUrl: "pikachu.jpg", aceSpec: false},
				"trainer": &mockCard{id: 2, name: "博士の研究", cardType: 2, imageUrl: "professor.jpg", aceSpec: false},
				"energy":  &mockCard{id: 3, name: "基本電気エネルギー", cardType: 3, imageUrl: "energy.jpg", aceSpec: false},
			},
			returnDeck: func() *domainDeck.Deck {
				// モックカードを作成
				mainCard := &mockCard{id: 1, name: "ピカチュウ", cardType: 1, imageUrl: "pikachu.jpg", aceSpec: false}
				subCard := &mockCard{id: 2, name: "博士の研究", cardType: 2, imageUrl: "professor.jpg", aceSpec: false}

				// モックデッキカードを作成
				card1 := &mockCard{id: 1, name: "ピカチュウ", cardType: 1, imageUrl: "pikachu.jpg", aceSpec: false}
				card2 := &mockCard{id: 2, name: "博士の研究", cardType: 2, imageUrl: "professor.jpg", aceSpec: false}
				card3 := &mockCard{id: 3, name: "基本電気エネルギー", cardType: 3, imageUrl: "energy.jpg", aceSpec: false}

				deckCard1 := domainDeck.NewDeckCard(card1, 4)
				deckCard2 := domainDeck.NewDeckCard(card2, 4)
				deckCard3 := domainDeck.NewDeckCard(card3, 52)

				deckCards := []domainDeck.DeckCard{*deckCard1, *deckCard2, *deckCard3}

				// デッキを作成
				deck, _ := domainDeck.NewDeck(1, "テストデッキ", "テスト用のデッキです", mainCard, subCard, deckCards)
				return deck
			}(),
			expectError: false,
			expectedDeckDto: &DeckDto{
				ID:          1,
				Name:        "テストデッキ",
				Description: "テスト用のデッキです",
				MainCard: &CardDto{
					ID:       1,
					Name:     "ピカチュウ",
					Category: "pokemon",
					ImageURL: "pikachu.jpg",
				},
				SubCard: &CardDto{
					ID:       2,
					Name:     "博士の研究",
					Category: "trainer",
					ImageURL: "professor.jpg",
				},
				Cards: []DeckCardWithQtyDto{
					{
						ID:       1,
						Name:     "ピカチュウ",
						Category: "pokemon",
						ImageURL: "pikachu.jpg",
						Quantity: 4,
					},
					{
						ID:       2,
						Name:     "博士の研究",
						Category: "trainer",
						ImageURL: "professor.jpg",
						Quantity: 4,
					},
					{
						ID:       3,
						Name:     "基本電気エネルギー",
						Category: "energy",
						ImageURL: "energy.jpg",
						Quantity: 52,
					},
				},
			},
		},
		"error_invalid_main_card_category": {
			request: &CreateDeckRequestDto{
				Name:        "テストデッキ",
				Description: "テスト用のデッキです",
				MainCardID: &CardIDDto{
					Id:       1,
					Category: "invalid",
				},
				Cards: []DeckCardRequestDto{},
			},
			mockCards:      map[string]domain.Card{},
			returnDeck:     nil,
			expectError:    true,
			expectedErrMsg: "invalid main card category",
		},
		"error_repository_failure": {
			request: &CreateDeckRequestDto{
				Name:        "テストデッキ",
				Description: "テスト用のデッキです",
				Cards: []DeckCardRequestDto{
					{
						Id:       1,
						Category: "pokemon",
						Quantity: 4,
					},
					{
						Id:       2,
						Category: "trainer",
						Quantity: 4,
					},
					{
						Id:       3,
						Category: "energy",
						Quantity: 52,
					},
				},
			},
			mockCards: map[string]domain.Card{
				"pokemon": &mockCard{id: 1, name: "ピカチュウ", cardType: 1, imageUrl: "pikachu.jpg", aceSpec: false},
				"trainer": &mockCard{id: 2, name: "博士の研究", cardType: 2, imageUrl: "professor.jpg", aceSpec: false},
				"energy":  &mockCard{id: 3, name: "基本電気エネルギー", cardType: 3, imageUrl: "energy.jpg", aceSpec: false},
			},
			returnDeck:     nil,
			expectError:    true,
			expectedErrMsg: "repository error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// モックのセットアップ
			mockDeckRepo := new(mockDeckRepository)
			mockCardRepo := new(mockCardRepository)

			// カードリポジトリのモック動作を設定
			for key, card := range tt.mockCards {
				parts := mock.AnythingOfType("domain.CardType")
				if key == "pokemon" {
					mockCardRepo.On("FindCardById", mock.Anything, 1, domain.Pokemon).Return(card, nil)
				} else if key == "trainer" {
					mockCardRepo.On("FindCardById", mock.Anything, 2, domain.Trainer).Return(card, nil)
				} else if key == "energy" {
					mockCardRepo.On("FindCardById", mock.Anything, 3, domain.Energy).Return(card, nil)
				} else {
					id := 0
					if key == "pokemon" {
						id = 1
					} else if key == "trainer" {
						id = 2
					} else if key == "energy" {
						id = 3
					}
					mockCardRepo.On("FindCardById", mock.Anything, id, parts).Return(card, nil)
				}
			}

			// エラーケースの設定
			if tt.expectError && tt.expectedErrMsg == "repository error" {
				// リポジトリエラーを模擬
				mockDeckRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
			} else if !tt.expectError {
				// 成功ケース
				mockDeckRepo.On("Create", mock.Anything, mock.Anything).Return(tt.returnDeck, nil)
			}

			// テスト対象のユースケースを作成
			useCase := NewCreateDeckUseCase(mockDeckRepo, mockCardRepo)

			// テスト実行
			result, err := useCase.Execute(context.Background(), tt.request)

			// アサーション
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				log.Println("err", err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedDeckDto.ID, result.ID)
				assert.Equal(t, tt.expectedDeckDto.Name, result.Name)
				assert.Equal(t, tt.expectedDeckDto.Description, result.Description)

				// メインカードの確認
				if tt.expectedDeckDto.MainCard != nil {
					assert.NotNil(t, result.MainCard)
					assert.Equal(t, tt.expectedDeckDto.MainCard.ID, result.MainCard.ID)
					assert.Equal(t, tt.expectedDeckDto.MainCard.Name, result.MainCard.Name)
					assert.Equal(t, tt.expectedDeckDto.MainCard.Category, result.MainCard.Category)
					assert.Equal(t, tt.expectedDeckDto.MainCard.ImageURL, result.MainCard.ImageURL)
				} else {
					assert.Nil(t, result.MainCard)
				}

				// サブカードの確認
				if tt.expectedDeckDto.SubCard != nil {
					assert.NotNil(t, result.SubCard)
					assert.Equal(t, tt.expectedDeckDto.SubCard.ID, result.SubCard.ID)
					assert.Equal(t, tt.expectedDeckDto.SubCard.Name, result.SubCard.Name)
					assert.Equal(t, tt.expectedDeckDto.SubCard.Category, result.SubCard.Category)
					assert.Equal(t, tt.expectedDeckDto.SubCard.ImageURL, result.SubCard.ImageURL)
				} else {
					assert.Nil(t, result.SubCard)
				}

				// カードリストの確認
				assert.Equal(t, len(tt.expectedDeckDto.Cards), len(result.Cards))
				for i, expectedCard := range tt.expectedDeckDto.Cards {
					assert.Equal(t, expectedCard.ID, result.Cards[i].ID)
					assert.Equal(t, expectedCard.Name, result.Cards[i].Name)
					assert.Equal(t, expectedCard.Category, result.Cards[i].Category)
					assert.Equal(t, expectedCard.ImageURL, result.Cards[i].ImageURL)
					assert.Equal(t, expectedCard.Quantity, result.Cards[i].Quantity)
				}
			}
		})
	}
}
