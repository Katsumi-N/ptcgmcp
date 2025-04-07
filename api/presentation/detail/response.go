package detail

type PokemonCardResponse struct {
	Result  bool        `json:"result"`
	Pokemon PokemonCard `json:"pokemon"`
}

type TrainerCardResponse struct {
	Result  bool        `json:"result"`
	Trainer TrainerCard `json:"trainer"`
}

type EnergyCardResponse struct {
	Result bool       `json:"result"`
	Energy EnergyCard `json:"energy"`
}

type PokemonCard struct {
	Id                 int             `json:"id"`
	Name               string          `json:"name"`
	EnergyType         string          `json:"energy_type"`
	Hp                 int             `json:"hp"`
	Ability            string          `json:"ability"`
	AbilityDescription string          `json:"ability_description"`
	ImageUrl           string          `json:"image_url"`
	Regulation         string          `json:"regulation"`
	Expansion          string          `json:"expansion"`
	Attacks            []PokemonAttack `json:"attacks"`
}

type PokemonAttack struct {
	Name           string `json:"name"`
	RequiredEnergy string `json:"required_energy"`
	Damage         string `json:"damage"`
	Description    string `json:"description"`
}

type TrainerCard struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	TrainerType string `json:"trainer_type"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}

type EnergyCard struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImageUrl    string `json:"image_url"`
	Description string `json:"description"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}
