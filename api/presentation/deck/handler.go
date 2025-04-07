package deck

import (
	deckUseCase "api/application/deck"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type deckHandler struct {
	listDeckUseCase     deckUseCase.IListDeckUseCase
	createDeckUseCase   deckUseCase.ICreateDeckUseCase
	validateDeckUseCase deckUseCase.IValidateDeckUseCase
	updateDeckUseCase   deckUseCase.IUpdateDeckUseCase
	deleteDeckUseCase   deckUseCase.IDeleteDeckUseCase
}

func NewDeckHandler(
	listDeckUseCase deckUseCase.IListDeckUseCase,
	createDeckUseCase deckUseCase.ICreateDeckUseCase,
	validateDeckUseCase deckUseCase.IValidateDeckUseCase,
	updateDeckUseCase deckUseCase.IUpdateDeckUseCase,
	deleteDeckUseCase deckUseCase.IDeleteDeckUseCase,
) *deckHandler {
	return &deckHandler{
		listDeckUseCase:     listDeckUseCase,
		createDeckUseCase:   createDeckUseCase,
		validateDeckUseCase: validateDeckUseCase,
		updateDeckUseCase:   updateDeckUseCase,
		deleteDeckUseCase:   deleteDeckUseCase,
	}
}

// GetUserDecks godoc
// @Summary Get user decks
// @Tags deck
// @Accept json
// @Produce json
// Removed @Security ApiKeyAuth
// @Success 200 {object} getUserDecksResponse
// @Router /v1/decks [get]
func (h *deckHandler) GetAllDecks(c echo.Context) error {
	// ユースケースを実行
	decks, err := h.listDeckUseCase.GetAllDecks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}

	// レスポンスを生成
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"decks":  decks,
	})
}

// CreateDeck godoc
// @Summary Create a new deck
// @Tags deck
// @Accept json
// @Produce json
// @Param request body createDeckRequest true "Deck information"
// @Success 200 {object} createDeckResponse
// @Router /v1/decks/create [post]
func (h *deckHandler) CreateDeck(c echo.Context) error {

	var req createDeckRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "Invalid request",
		})
	}

	// ユースケース用のDTOを作成
	requestDto := &deckUseCase.CreateDeckRequestDto{
		Name:        req.Name,
		Description: req.Description,
		Cards:       make([]deckUseCase.DeckCardRequestDto, 0, len(req.Cards)),
	}

	// メインカードとサブカードがある場合は設定
	if req.MainCard != nil {
		requestDto.MainCardID = &deckUseCase.CardIDDto{
			Id:       req.MainCard.ID,
			Category: req.MainCard.Category,
		}
	}

	if req.SubCard != nil {
		requestDto.SubCardID = &deckUseCase.CardIDDto{
			Id:       req.SubCard.ID,
			Category: req.SubCard.Category,
		}
	}

	// カードリストの設定
	for _, card := range req.Cards {
		requestDto.Cards = append(requestDto.Cards, deckUseCase.DeckCardRequestDto{
			Id:       card.ID,
			Category: card.Category,
			Quantity: card.Quantity,
		})
	}

	// ユースケースを実行
	deck, err := h.createDeckUseCase.Execute(c.Request().Context(), requestDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}

	// レスポンスを生成
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"deck":   deck,
	})
}

// ValidateDeck godoc
// @Summary Validate a deck
// @Tags deck
// @Accept json
// @Produce json
// @Param request body validateDeckRequest true "Deck information"
// @Success 200 {object} validateDeckResponse
// @Router /v1/decks/validate [post]
func (h *deckHandler) ValidateDeck(c echo.Context) error {

	// リクエストをバインド
	var req validateDeckRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "Invalid request",
		})
	}

	// ユースケース用のDTOを作成
	requestDto := &deckUseCase.ValidateDeckRequestDto{
		Name:        req.Name,
		Description: req.Description,
		Cards:       make([]deckUseCase.DeckCardRequestDto, 0, len(req.Cards)),
	}

	// メインカードとサブカードがある場合は設定
	if req.MainCard != nil {
		requestDto.MainCardID = &deckUseCase.CardIDDto{
			Id:       req.MainCard.ID,
			Category: req.MainCard.Category,
		}
	}

	if req.SubCard != nil {
		requestDto.SubCardID = &deckUseCase.CardIDDto{
			Id:       req.SubCard.ID,
			Category: req.SubCard.Category,
		}
	}

	// カードリストの設定
	for _, card := range req.Cards {
		requestDto.Cards = append(requestDto.Cards, deckUseCase.DeckCardRequestDto{
			Id:       card.ID,
			Category: card.Category,
			Quantity: card.Quantity,
		})
	}

	// ユースケースを実行
	result, err := h.validateDeckUseCase.Execute(c.Request().Context(), requestDto)
	if err != nil {
		// エラーメッセージからHTTPステータスを判断
		if errors.Is(err, errors.New("invalid main card category")) ||
			errors.Is(err, errors.New("invalid sub card category")) ||
			errors.Is(err, errors.New("invalid card category")) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"result": false,
				"error":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}

	// レスポンスを生成
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result":   true,
		"is_valid": result.IsValid,
		"errors":   result.Errors,
	})
}

// UpdateDeck godoc
// @Summary Update a deck
// @Tags deck
// @Accept json
// @Produce json
// @Param id path int true "Deck ID"
// @Param request body updateDeckRequest true "Deck information"
// @Success 200 {object} updateDeckResponse
// @Router /v1/decks/edit/{id} [post]
func (h *deckHandler) UpdateDeck(c echo.Context) error {

	// デッキIDをパスパラメータから取得
	deckIdStr := c.Param("id")
	deckId, err := strconv.Atoi(deckIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "不正なデッキIDです",
		})
	}

	// リクエストをバインド
	var req updateDeckRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "リクエスト形式が不正です",
		})
	}

	// ユースケース用のDTOを作成
	requestDto := &deckUseCase.UpdateDeckRequestDto{
		Name:        req.Name,
		Description: req.Description,
		Cards:       make([]deckUseCase.DeckCardRequestDto, 0, len(req.Cards)),
	}

	// メインカードとサブカードがある場合は設定
	if req.MainCard != nil {
		requestDto.MainCardID = &deckUseCase.CardIDDto{
			Id:       req.MainCard.ID,
			Category: req.MainCard.Category,
		}
	}

	if req.SubCard != nil {
		requestDto.SubCardID = &deckUseCase.CardIDDto{
			Id:       req.SubCard.ID,
			Category: req.SubCard.Category,
		}
	}

	// カードリストの設定
	for _, card := range req.Cards {
		requestDto.Cards = append(requestDto.Cards, deckUseCase.DeckCardRequestDto{
			Id:       card.ID,
			Category: card.Category,
			Quantity: card.Quantity,
		})
	}

	// ユースケースを実行
	deck, err := h.updateDeckUseCase.Execute(c.Request().Context(), deckId, requestDto)
	if err != nil {
		if errors.Is(err, errors.New("デッキが見つかりません")) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"result": false,
				"error":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}

	// レスポンスを生成
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"deck":   deck,
	})
}

// DeleteDeck godoc
// @Summary Delete a deck
// @Tags deck
// @Accept json
// @Produce json
// @Param id path int true "Deck ID"
// @Success 200 {object} deleteDeckResponse
// @Router /v1/decks/delete/{id} [delete]
func (h *deckHandler) DeleteDeck(c echo.Context) error {

	deckIdStr := c.Param("id")
	deckId, err := strconv.Atoi(deckIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "不正なデッキIDです",
		})
	}

	err = h.deleteDeckUseCase.DeleteDeck(c.Request().Context(), deckId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result":  true,
		"message": "デッキが削除されました",
	})
}

// GetDeckById godoc
// @Summary Get deck by ID
// @Tags deck
// @Accept json
// @Produce json
// @Param id path int true "Deck ID"
// @Success 200 {object} getDeckByIdResponse
// @Router /v1/decks/detail/{id} [get]
func (h *deckHandler) GetDeckById(c echo.Context) error {
	deckIdStr := c.Param("id")
	deckId, err := strconv.Atoi(deckIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": false,
			"error":  "不正なデッキIDです",
		})
	}

	deck, err := h.listDeckUseCase.GetDeckById(c.Request().Context(), deckId)
	if err != nil {
		if errors.Is(err, errors.New("デッキが見つかりません")) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"result": false,
				"error":  "デッキが見つかりません",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"result": false,
			"error":  err.Error(),
		})
	}

	// レスポンスを生成
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"deck":   deck,
	})
}
