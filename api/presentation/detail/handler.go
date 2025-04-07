package detail

import (
	"api/application/detail"
	"api/domain"
	errDomain "api/domain/error"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type detailHandler struct {
	detailUseCase *detail.FetchDetailUseCase
}

func NewDetailHandler(detailUseCase *detail.FetchDetailUseCase) *detailHandler {
	return &detailHandler{
		detailUseCase: detailUseCase,
	}
}

// FetchDetail godoc
// @Summary fetch pokemon/trainer/energy detail
// @Tags detail
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param card_type path string true "card_type"
// @Success 200 {object} getDetailResponse
// @Router /v1/cards/detail/{card_type}/{id} [get]
func (h *detailHandler) FetchDetail(c echo.Context) error {
	cardType := c.Param("card_type")
	id := c.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(500, "id must be integer")
	}

	switch domain.StringToCardType[cardType] {
	case domain.Pokemon:
		pokemon, err := h.detailUseCase.FetchPokemonDetail(c.Request().Context(), iid)
		if err != nil {
			if err == errDomain.NotFoundErr {
				return c.JSON(404, fmt.Sprintf("pokemon id %d not found", iid))
			}
			return c.JSON(500, err.Error())
		}

		attacks := lo.Map(pokemon.Attacks, func(attack detail.PokemonAttack, _ int) PokemonAttack {
			return PokemonAttack(attack)
		})

		return c.JSON(200, PokemonCardResponse{
			Result: true,
			Pokemon: PokemonCard{
				Id:                 pokemon.Id,
				Name:               pokemon.Name,
				EnergyType:         pokemon.EnergyType,
				Hp:                 pokemon.Hp,
				Ability:            pokemon.Ability,
				AbilityDescription: pokemon.AbilityDescription,
				ImageUrl:           pokemon.ImageUrl,
				Regulation:         pokemon.Regulation,
				Expansion:          pokemon.Expansion,
				Attacks:            attacks,
			},
		})
	case domain.Trainer:
		trainer, err := h.detailUseCase.FetchTrainerDetail(c.Request().Context(), iid)
		if err != nil {
			return c.JSON(500, err.Error())
		}

		return c.JSON(200, TrainerCardResponse{
			Result: true,
			Trainer: TrainerCard{
				Id:          trainer.Id,
				Name:        trainer.Name,
				TrainerType: trainer.TrainerType,
				Description: trainer.Description,
				ImageUrl:    trainer.ImageUrl,
				Regulation:  trainer.Regulation,
				Expansion:   trainer.Expansion,
			},
		})

	case domain.Energy:
		energy, err := h.detailUseCase.FetchEnergyDetail(c.Request().Context(), iid)
		if err != nil {
			return c.JSON(500, err.Error())
		}

		return c.JSON(200, EnergyCardResponse{
			Result: true,
			Energy: EnergyCard{
				Id:          energy.Id,
				Name:        energy.Name,
				ImageUrl:    energy.ImageUrl,
				Description: energy.Description,
				Regulation:  energy.Regulation,
				Expansion:   energy.Expansion,
			},
		})

	default:
		return c.JSON(500, "invalid card type")
	}

}
