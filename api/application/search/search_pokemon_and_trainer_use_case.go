package search

import (
	"api/application/search/energy"
	"api/application/search/pokemon"
	"api/application/search/trainer"
	"context"
	"fmt"

	"github.com/samber/lo"
)

type SearchPokemonAndTrainerUseCase struct {
	pokemonQueryService pokemon.PokemonQueryService
	trainerQueryService trainer.TrainerQueryService
	energyQueryService  energy.EnergyQueryService
}

func NewSearchPokemonAndTrainerUseCase(
	pokemonQueryService pokemon.PokemonQueryService,
	trainerQueryService trainer.TrainerQueryService,
	energyQueryService energy.EnergyQueryService,
) *SearchPokemonAndTrainerUseCase {
	return &SearchPokemonAndTrainerUseCase{
		pokemonQueryService: pokemonQueryService,
		trainerQueryService: trainerQueryService,
		energyQueryService:  energyQueryService,
	}
}

type SearchPokemonAndTrainerUseCaseDto struct {
	Pokemons []*pokemon.SearchPokemonUseCaseDto `json:"pokemons"`
	Trainers []*trainer.SearchTrainerUseCaseDto `json:"trainers"`
	Energies []*energy.SearchEnergyUseCaseDto   `json:"energies"`
}

func (uc *SearchPokemonAndTrainerUseCase) SearchPokemonAndTrainerList(ctx context.Context, q string) (*SearchPokemonAndTrainerUseCaseDto, error) {
	searchPokemonList, err := uc.pokemonQueryService.SearchPokemonList(ctx, q)
	if err != nil {
		return nil, err
	}

	searchTrainerList, err := uc.trainerQueryService.SearchTrainerList(ctx, q)
	if err != nil {
		return nil, err
	}

	searchEnergyList, err := uc.energyQueryService.SearchEnergyList(ctx, q)
	if err != nil {
		return nil, err
	}

	pokemons := lo.Map(searchPokemonList, func(f *pokemon.SearchPokemonList, _ int) *pokemon.SearchPokemonUseCaseDto {
		return &pokemon.SearchPokemonUseCaseDto{
			ID:         fmt.Sprintf("%v", f.ID),
			Name:       f.Name,
			EnergyType: f.EnergyType,
			Hp:         f.Hp,
			ImageURL:   f.ImageURL,
		}
	})

	trainers := lo.Map(searchTrainerList, func(f *trainer.SearchTrainerList, _ int) *trainer.SearchTrainerUseCaseDto {
		return &trainer.SearchTrainerUseCaseDto{
			ID:          fmt.Sprintf("%v", f.ID),
			Name:        f.Name,
			TrainerType: f.TrainerType,
			ImageURL:    f.ImageURL,
		}
	})

	energies := lo.Map(searchEnergyList, func(f *energy.SearchEnergyList, _ int) *energy.SearchEnergyUseCaseDto {
		return &energy.SearchEnergyUseCaseDto{
			ID:          fmt.Sprintf("%v", f.ID),
			Name:        f.Name,
			ImageURL:    f.ImageURL,
			Description: f.Description,
		}
	})

	dto := &SearchPokemonAndTrainerUseCaseDto{
		Pokemons: pokemons,
		Trainers: trainers,
		Energies: energies,
	}

	return dto, nil
}

func (uc *SearchPokemonAndTrainerUseCase) SearchPokemonList(ctx context.Context, q string) (*SearchPokemonAndTrainerUseCaseDto, error) {
	searchPokemonList, err := uc.pokemonQueryService.SearchPokemonList(ctx, q)
	if err != nil {
		return nil, err
	}

	pokemons := lo.Map(searchPokemonList, func(f *pokemon.SearchPokemonList, _ int) *pokemon.SearchPokemonUseCaseDto {
		return &pokemon.SearchPokemonUseCaseDto{
			ID:         fmt.Sprintf("%v", f.ID),
			Name:       f.Name,
			EnergyType: f.EnergyType,
			Hp:         f.Hp,
			ImageURL:   f.ImageURL,
		}
	})

	dto := &SearchPokemonAndTrainerUseCaseDto{
		Pokemons: pokemons,
		Trainers: make([]*trainer.SearchTrainerUseCaseDto, 0),
		Energies: make([]*energy.SearchEnergyUseCaseDto, 0),
	}

	return dto, nil
}

func (uc *SearchPokemonAndTrainerUseCase) SearchTrainerList(ctx context.Context, q string) (*SearchPokemonAndTrainerUseCaseDto, error) {
	searchTrainerList, err := uc.trainerQueryService.SearchTrainerList(ctx, q)
	if err != nil {
		return nil, err
	}

	trainers := lo.Map(searchTrainerList, func(f *trainer.SearchTrainerList, _ int) *trainer.SearchTrainerUseCaseDto {
		return &trainer.SearchTrainerUseCaseDto{
			ID:          fmt.Sprintf("%v", f.ID),
			Name:        f.Name,
			TrainerType: f.TrainerType,
			ImageURL:    f.ImageURL,
		}
	})

	dto := &SearchPokemonAndTrainerUseCaseDto{
		Pokemons: make([]*pokemon.SearchPokemonUseCaseDto, 0),
		Trainers: trainers,
		Energies: make([]*energy.SearchEnergyUseCaseDto, 0),
	}

	return dto, nil
}

func (uc *SearchPokemonAndTrainerUseCase) SearchEnergyList(ctx context.Context, q string) (*SearchPokemonAndTrainerUseCaseDto, error) {
	searchEnergyList, err := uc.energyQueryService.SearchEnergyList(ctx, q)
	if err != nil {
		return nil, err
	}

	energies := lo.Map(searchEnergyList, func(f *energy.SearchEnergyList, _ int) *energy.SearchEnergyUseCaseDto {
		return &energy.SearchEnergyUseCaseDto{
			ID:          fmt.Sprintf("%v", f.ID),
			Name:        f.Name,
			ImageURL:    f.ImageURL,
			Description: f.Description,
		}
	})

	dto := &SearchPokemonAndTrainerUseCaseDto{
		Pokemons: make([]*pokemon.SearchPokemonUseCaseDto, 0),
		Trainers: make([]*trainer.SearchTrainerUseCaseDto, 0),
		Energies: energies,
	}

	return dto, nil
}
