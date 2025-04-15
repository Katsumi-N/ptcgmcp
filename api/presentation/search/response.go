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
	TrainerType string `json:"trainer_type"`
	ImageURL    string `json:"image_url"`
}

type energy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type searchDeckResponse struct {
	Result bool            `json:"result"`
	Decks  []*searchedDeck `json:"decks"`
}

type searchedDeck struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	MainCard    *deckCard   `json:"main_card"`
	SubCard     *deckCard   `json:"sub_card"`
	Cards       []*deckCard `json:"cards"`
}
type deckCard struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	ImageURL string `json:"image_url"`
}
