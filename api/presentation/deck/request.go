package deck

// CreateDeck Request
type createDeckRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	MainCard    *cardIDRequest    `json:"main_card,omitempty"`
	SubCard     *cardIDRequest    `json:"sub_card,omitempty"`
	Cards       []deckCardRequest `json:"cards"`
}

type cardIDRequest struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
}

type deckCardRequest struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
}

// ValidateDeck Request
type validateDeckRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	MainCard    *cardIDRequest    `json:"main_card,omitempty"`
	SubCard     *cardIDRequest    `json:"sub_card,omitempty"`
	Cards       []deckCardRequest `json:"cards"`
}

type updateDeckRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	MainCard    *cardIDRequest    `json:"main_card,omitempty"`
	SubCard     *cardIDRequest    `json:"sub_card,omitempty"`
	Cards       []deckCardRequest `json:"cards"`
}
