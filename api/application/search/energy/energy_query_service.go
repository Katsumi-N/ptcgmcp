package energy

import "context"

type SearchEnergyList struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

type EnergyQueryService interface {
	SearchEnergyList(ctx context.Context, q string) ([]*SearchEnergyList, error)
}
