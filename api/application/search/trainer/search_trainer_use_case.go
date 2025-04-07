package trainer

import (
	"context"
	"fmt"
)

type SearchTrainerUseCase struct {
	trainerQueryService TrainerQueryService
}

type SearchTrainerUseCaseDto struct {
	ID          string
	Name        string
	TrainerType string
	ImageURL    string
}

func (uc *SearchTrainerUseCase) SearchTrainerList(ctx context.Context, q string) ([]*SearchTrainerUseCaseDto, error) {
	searchTrainerList, err := uc.trainerQueryService.SearchTrainerList(ctx, q)
	if err != nil {
		return nil, err
	}

	var dtoList []*SearchTrainerUseCaseDto
	for _, f := range searchTrainerList {
		dto := &SearchTrainerUseCaseDto{
			ID:          fmt.Sprintf("%v", f.ID),
			Name:        f.Name,
			TrainerType: f.TrainerType,
			ImageURL:    f.ImageURL,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}
