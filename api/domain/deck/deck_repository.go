package deck

import (
	"api/domain"
	"context"
)

type DeckRepository interface {
	// デッキの作成
	Create(ctx context.Context, deck *Deck) (*Deck, error)

	FindAll(ctx context.Context) ([]*Deck, error)

	// デッキの詳細取得
	FindById(ctx context.Context, id int) (*Deck, error)

	// デッキの更新
	Update(ctx context.Context, deck *Deck) error

	// デッキの削除
	Delete(ctx context.Context, id int) error
}

// カード情報を取得するためのリポジトリ
type CardRepository interface {
	// カードIDとタイプからカード情報を取得
	FindCardById(ctx context.Context, cardId int, cardType domain.CardType) (domain.Card, error)
}
