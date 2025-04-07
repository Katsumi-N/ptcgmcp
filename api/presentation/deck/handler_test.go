package deck

import (
	deckUseCase "api/application/deck"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 認証ミドルウェアのモック
type mockAuthMiddleware struct{}

func (m *mockAuthMiddleware) GetAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTトークンを模倣
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": "test-user-id",
			})
			// コンテキストにJWTトークンをセット
			c.Set("user", token)
			return next(c)
		}
	}
}

// モックユースケース - CreateDeck
type mockCreateDeckUseCase struct {
	mock.Mock
}

func (m *mockCreateDeckUseCase) Execute(ctx context.Context, request *deckUseCase.CreateDeckRequestDto) (*deckUseCase.DeckDto, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*deckUseCase.DeckDto), args.Error(1)
}

// モックユースケース - ListDeck
type mockListDeckUseCase struct {
	mock.Mock
}

func (m *mockListDeckUseCase) GetAllDecks(ctx context.Context) ([]*deckUseCase.DeckDto, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*deckUseCase.DeckDto), args.Error(1)
}

func (m *mockListDeckUseCase) GetDeckById(ctx context.Context, id int) (*deckUseCase.DeckDto, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*deckUseCase.DeckDto), args.Error(1)
}

// モックユースケース - ValidateDeck
type mockValidateDeckUseCase struct {
	mock.Mock
}

func (m *mockValidateDeckUseCase) Execute(ctx context.Context, request *deckUseCase.ValidateDeckRequestDto) (*deckUseCase.ValidateDeckResponseDto, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*deckUseCase.ValidateDeckResponseDto), args.Error(1)
}

type mockUpdateDeckUseCase struct {
	mock.Mock
}

func (m *mockUpdateDeckUseCase) Execute(ctx context.Context, id int, request *deckUseCase.UpdateDeckRequestDto) (*deckUseCase.DeckDto, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*deckUseCase.DeckDto), args.Error(1)
}

type mockDeleteDeckUseCase struct {
	mock.Mock
}

func (m *mockDeleteDeckUseCase) DeleteDeck(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func TestCreateDeck(t *testing.T) {
	// テストケースの定義
	tests := map[string]struct {
		requestBody        string
		mockReturn         *deckUseCase.DeckDto
		mockError          error
		expectedStatusCode int
		expectedResult     bool
	}{
		"success": {
			requestBody: `{
				"name": "テストデッキ",
				"description": "テスト用のデッキです",
				"main_card": {
					"id": 1,
					"category": "pokemon"
				},
				"sub_card": {
					"id": 2,
					"category": "trainer"
				},
				"cards": [
					{
						"id": 1,
						"category": "pokemon",
						"quantity": 4
					},
					{
						"id": 2,
						"category": "trainer",
						"quantity": 4
					},
					{
						"id": 3,
						"category": "energy",
						"quantity": 52
					}
				]
			}`,
			mockReturn: &deckUseCase.DeckDto{
				ID:          1,
				Name:        "テストデッキ",
				Description: "テスト用のデッキです",
				MainCard: &deckUseCase.CardDto{
					ID:       1,
					Name:     "ピカチュウ",
					Category: "pokemon",
					ImageURL: "pikachu.jpg",
				},
				SubCard: &deckUseCase.CardDto{
					ID:       2,
					Name:     "博士の研究",
					Category: "trainer",
					ImageURL: "professor.jpg",
				},
				Cards: []deckUseCase.DeckCardWithQtyDto{
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
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
			expectedResult:     true,
		},
		"invalid_request": {
			requestBody:        `{"name": ""}`,
			mockReturn:         nil,
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     false,
		},
		"server_error": {
			requestBody: `{
				"name": "テストデッキ",
				"description": "テスト用のデッキです",
				"cards": []
			}`,
			mockReturn:         nil,
			mockError:          errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResult:     false,
		},
		"invalid_main_card_category": {
			requestBody: `{
				"name": "テストデッキ",
				"description": "テスト用のデッキです",
				"main_card": {
					"id": 1,
					"category": "invalid"
				},
				"cards": []
			}`,
			mockReturn:         nil,
			mockError:          errors.New("invalid main card category"),
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Echoのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/v1/decks/create", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 認証ミドルウェアを適用
			authMiddleware := &mockAuthMiddleware{}
			middleware := authMiddleware.GetAuthMiddleware()
			h := middleware(func(c echo.Context) error {
				// モックの作成
				mockCreateDeckUC := new(mockCreateDeckUseCase)
				mockListDeckUC := new(mockListDeckUseCase)
				mockValidateDeckUC := new(mockValidateDeckUseCase)
				mockUpdateDeckUseCase := new(mockUpdateDeckUseCase)
				mockDeleteDeckUseCase := new(mockDeleteDeckUseCase)

				// モックの振る舞いを設定（正しいパッケージパスとジェネリックな引数を指定）
				if tt.requestBody != `{"name": ""}` { // 無効なリクエストの場合はモックは呼び出されない
					mockCreateDeckUC.On("Execute", mock.Anything, mock.AnythingOfType("*deck.CreateDeckRequestDto")).Return(tt.mockReturn, tt.mockError)
				}

				// ハンドラーの作成
				handler := NewDeckHandler(mockListDeckUC, mockCreateDeckUC, mockValidateDeckUC, mockUpdateDeckUseCase, mockDeleteDeckUseCase)

				// テスト対象の関数を呼び出し
				return handler.CreateDeck(c)
			})

			// テスト実行
			err := h(c)
			if err != nil {
				// 認証エラー（Unauthorized）以外の予期しないエラーの場合はテスト失敗
				t.Fatalf("Unexpected error: %v", err)
			}

			// アサーション
			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			// レスポンスのJSONをパース
			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Error parsing response JSON: %v", err)
			}

			// 結果の確認
			assert.Equal(t, tt.expectedResult, response["result"])

			// 成功した場合、デッキの内容も確認
			if tt.expectedResult {
				assert.NotNil(t, response["deck"])
			} else {
				// エラーの場合はエラーメッセージが含まれる
				if tt.expectedStatusCode != http.StatusBadRequest {
					assert.Contains(t, response, "error")
				}
			}
		})
	}
}
