package search

import (
	card "api/application/search"
	deck "api/application/search/deck"
	"api/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type searchHandler struct {
	searchCardUseCase *card.SearchPokemonAndTrainerUseCase
	SearchDeckUseCase *deck.SearchDeckUseCase
}

func NewSearchHandler(searchCardUseCase *card.SearchPokemonAndTrainerUseCase, searchDeckUseCase *deck.SearchDeckUseCase) searchHandler {
	return searchHandler{
		searchCardUseCase: searchCardUseCase,
		SearchDeckUseCase: searchDeckUseCase,
	}
}

// SearchCardList godoc
// @Summary Search card list
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {object} getProductsResponse
// @Router /v1/cards/search [get]
func (h *searchHandler) SearchCardList(c echo.Context) error {
	q := c.QueryParam("q")
	cardType := c.QueryParam("card_type")
	dto, err := func(cardType string) (*card.SearchPokemonAndTrainerUseCaseDto, error) {
		switch domain.StringToCardType[cardType] {
		case domain.Pokemon:
			return h.searchCardUseCase.SearchPokemonList(c.Request().Context(), q)
		case domain.Trainer:
			return h.searchCardUseCase.SearchTrainerList(c.Request().Context(), q)
		case domain.Energy:
			return h.searchCardUseCase.SearchEnergyList(c.Request().Context(), q)
		default:
			return h.searchCardUseCase.SearchPokemonAndTrainerList(c.Request().Context(), q)
		}
	}(cardType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var res searchCardResponse
	res.Result = true
	for _, dtoPokemon := range dto.Pokemons {
		res.Pokemons = append(res.Pokemons, &pokemon{
			ID:         dtoPokemon.ID,
			Name:       dtoPokemon.Name,
			EnergyType: dtoPokemon.EnergyType,
			Hp:         dtoPokemon.Hp,
			ImageURL:   dtoPokemon.ImageURL,
		})
	}

	for _, dtoTrainer := range dto.Trainers {
		res.Trainers = append(res.Trainers, &trainer{
			ID:          dtoTrainer.ID,
			Name:        dtoTrainer.Name,
			TrainerType: dtoTrainer.TrainerType,
			ImageURL:    dtoTrainer.ImageURL,
		})
	}

	for _, dtoEnergy := range dto.Energies {
		res.Energies = append(res.Energies, &energy{
			ID:          dtoEnergy.ID,
			Name:        dtoEnergy.Name,
			Description: dtoEnergy.Description,
			ImageURL:    dtoEnergy.ImageURL,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// SearchDeckList godoc
// @Summary Search deck list
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {object} getDecksResponse
// @Router /v1/decks/search [get]
func (h *searchHandler) SearchDeckList(c echo.Context) error {
	q := c.QueryParam("q")
	dto, err := h.SearchDeckUseCase.SearchDeckList(c.Request().Context(), q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var res searchDeckResponse
	res.Result = true

	resultDecks := lo.Map(dto, func(f *deck.SearchDeckUseCaseDto, _ int) *searchedDeck {
		return &searchedDeck{
			ID:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			MainCard: &deckCard{
				ID:       f.MainCard.Id,
				Name:     f.MainCard.Name,
				Category: f.MainCard.Category,
				ImageURL: f.MainCard.ImageURL,
			},
			SubCard: &deckCard{
				ID:       f.SubCard.Id,
				Name:     f.SubCard.Name,
				Category: f.SubCard.Category,
				ImageURL: f.SubCard.ImageURL,
			},
			Cards: lo.Map(f.Cards, func(card deck.SearchDeckCardUseCaseDto, _ int) *deckCard {
				return &deckCard{
					ID:       card.Id,
					Name:     card.Name,
					Category: card.Category,
					Quantity: card.Quantity,
					ImageURL: card.ImageURL,
				}
			}),
		}
	})
	res.Decks = resultDecks
	return c.JSON(http.StatusOK, res)
}
