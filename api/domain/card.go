package domain

type Card interface {
	GetId() int
	GetName() string
	GetCardType() int
	GetImageUrl() string
	IsAceSpec() bool
}
