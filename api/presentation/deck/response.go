package deck

// GetUserDecks Response
type getUserDecksResponse struct {
	Result bool        `json:"result"`
	Decks  interface{} `json:"decks"`
}

// CreateDeck Response
type createDeckResponse struct {
	Result bool        `json:"result"`
	Deck   interface{} `json:"deck"`
}

// ValidateDeck Response
type validateDeckResponse struct {
	Result  bool     `json:"result"`
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

// UpdateDeck Response
type updateDeckResponse struct {
	Result bool        `json:"result"`
	Deck   interface{} `json:"deck,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// GetDeckById Response
type getDeckByIdResponse struct {
	Result bool        `json:"result"`
	Deck   interface{} `json:"deck,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// DeleteDeck Response
type deleteDeckResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
