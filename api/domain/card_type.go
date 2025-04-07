package domain

type CardType int

const (
	Pokemon CardType = iota + 1
	Trainer
	Energy
)

var CardTypeToString = map[CardType]string{
	Pokemon: "pokemon",
	Trainer: "trainer",
	Energy:  "energy",
}

var StringToCardType = map[string]CardType{
	"pokemon": Pokemon,
	"trainer": Trainer,
	"energy":  Energy,
}
