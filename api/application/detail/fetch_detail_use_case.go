package detail

import "context"

type FetchDetailUseCase struct {
	detailQueryService DetailQueryService
}

func NewFetchDetailUseCase(detailQueryService DetailQueryService) *FetchDetailUseCase {
	return &FetchDetailUseCase{
		detailQueryService: detailQueryService,
	}
}

func (uc *FetchDetailUseCase) FetchPokemonDetail(ctx context.Context, pokemonId int) (*Pokemon, error) {
	pokemon, err := uc.detailQueryService.FindPokemonDetail(ctx, pokemonId)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

func (uc *FetchDetailUseCase) FetchTrainerDetail(ctx context.Context, trainerId int) (*Trainer, error) {
	trainer, err := uc.detailQueryService.FindTrainerDetail(ctx, trainerId)
	if err != nil {
		return nil, err
	}

	return trainer, nil
}

func (uc *FetchDetailUseCase) FetchEnergyDetail(ctx context.Context, energyId int) (*Energy, error) {
	energy, err := uc.detailQueryService.FindEnergyDetail(ctx, energyId)
	if err != nil {
		return nil, err
	}

	return energy, nil
}
