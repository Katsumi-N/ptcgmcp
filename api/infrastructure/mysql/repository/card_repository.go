package repository

import (
	"api/domain"
	"api/domain/deck"
	"api/domain/energy"
	domainErr "api/domain/error"
	"api/domain/pokemon"
	"api/domain/trainer"
	"api/infrastructure/mysql/db"
	"context"
	"database/sql"
	"errors"
)

type cardRepository struct{}

func NewCardRepository() deck.CardRepository {
	return &cardRepository{}
}

// カードIDとタイプからカード情報を取得
func (r *cardRepository) FindCardById(ctx context.Context, cardId int, cardType domain.CardType) (domain.Card, error) {
	query := db.GetQuery(ctx)
	switch cardType {
	case domain.Pokemon:
		pokemonRow, err := query.PokemonFindById(ctx, int64(cardId))
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domainErr.NotFoundErr
			}
			return nil, err
		}

		p, err := pokemon.NewPokemon(
			int(pokemonRow.ID),
			pokemonRow.Name,
			pokemonRow.EnergyType,
			int(pokemonRow.Hp),
			pokemonRow.Ability.String,
			pokemonRow.AbilityDescription.String,
			pokemonRow.ImageUrl,
			pokemonRow.Regulation,
			pokemonRow.Expansion,
			nil, // デッキ情報としてワザは不要
		)
		if err != nil {
			return nil, err
		}

		return p, nil

	case domain.Trainer:
		// エネルギーカードとトレーナーカードが逆になっていたため修正
		trainerRow, err := query.TrainerFindById(ctx, int64(cardId))
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domainErr.NotFoundErr
			}
			return nil, err
		}

		t, err := trainer.NewTrainer(
			int(trainerRow.ID),
			trainerRow.Name,
			trainerRow.TrainerType,
			trainerRow.Description,
			trainerRow.ImageUrl,
			trainerRow.Regulation,
			trainerRow.Expansion,
		)

		return t, err
	case domain.Energy:
		// エネルギーカードとトレーナーカードが逆になっていたため修正
		energyCard, err := query.EnergyFindById(ctx, int64(cardId))
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domainErr.NotFoundErr
			}
			return nil, err
		}
		e, err := energy.NewEnergy(
			int(energyCard.ID),
			energyCard.Name,
			energyCard.ImageUrl,
			energyCard.Regulation,
			energyCard.Expansion,
		)

		return e, err
	default:
		return nil, errors.New("invalid card type")
	}
}
