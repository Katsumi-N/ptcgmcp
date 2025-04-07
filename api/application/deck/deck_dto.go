package deck

type DeckDto struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	MainCard    *CardDto             `json:"main_card,omitempty"`
	SubCard     *CardDto             `json:"sub_card,omitempty"`
	Cards       []DeckCardWithQtyDto `json:"cards"`
}

type CardDto struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
}

type DeckCardWithQtyDto struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
	Quantity int    `json:"quantity"`
}
