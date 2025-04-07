package search

type searchCardResponse struct {
	Result   bool       `json:"result"`
	Pokemons []*pokemon `json:"pokemons"`
	Trainers []*trainer `json:"trainers"`
	Energies []*energy  `json:"energies"`
}

type pokemon struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Hp         int    `json:"hp"`
	EnergyType string `json:"energy_type"`
	ImageURL   string `json:"image_url"`
}

type trainer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TrainerType string `json:"energy_type"`
	ImageURL    string `json:"image_url"`
}

type energy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
