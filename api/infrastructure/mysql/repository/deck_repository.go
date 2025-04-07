package repository

import (
	"api/domain"
	"api/domain/deck"
	"api/infrastructure/mysql/db"
	"api/infrastructure/mysql/db/dbgen"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type deckRepository struct {
	cardRepository deck.CardRepository
}

// DeckRepositoryインターフェースの実装
func NewDeckRepository() deck.DeckRepository {
	return &deckRepository{
		cardRepository: NewCardRepository(),
	}
}

// デッキの作成
func (r *deckRepository) Create(ctx context.Context, d *deck.Deck) (*deck.Deck, error) {
	// トランザクション開始
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback() // 明示的にコミットされなければロールバックする

	qtx := dbgen.New(tx)

	// メインカードとサブカードのID、タイプを取得
	var mainCardID sql.NullInt64
	var mainCardTypeID sql.NullInt64
	var subCardID sql.NullInt64
	var subCardTypeID sql.NullInt64

	if d.GetMainCard() != nil {
		mainCardID.Int64 = int64(d.GetMainCard().GetId())
		mainCardID.Valid = true
		mainCardTypeID.Int64 = int64(d.GetMainCard().GetCardType())
		mainCardTypeID.Valid = true
	}

	if d.GetSubCard() != nil {
		subCardID.Int64 = int64(d.GetSubCard().GetId())
		subCardID.Valid = true
		subCardTypeID.Int64 = int64(d.GetSubCard().GetCardType())
		subCardTypeID.Valid = true
	}

	// デッキを作成
	deckRow, err := qtx.CreateDeck(ctx, dbgen.CreateDeckParams{
		Name:           d.GetName(),
		Description:    sql.NullString{String: d.GetDescription(), Valid: d.GetDescription() != ""},
		MainCardID:     mainCardID,
		MainCardTypeID: mainCardTypeID,
		SubCardID:      subCardID,
		SubCardTypeID:  subCardTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("デッキ作成エラー: %w", err)
	}

	// 作成されたデッキIDを取得
	insertedId, _ := deckRow.LastInsertId()

	// デッキカードを追加
	for _, card := range d.GetCards() {
		_, err = qtx.CreateDeckCard(ctx, dbgen.CreateDeckCardParams{
			DeckID:     insertedId, // ここにデッキIDを追加
			CardID:     int64(card.GetCard().GetId()),
			CardTypeID: int64(card.GetCard().GetCardType()),
			Quantity:   int32(card.GetQuantity()),
		})
		if err != nil {
			return nil, fmt.Errorf("デッキカード作成エラー: %w", err)
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	// 作成したデッキを取得して返す
	return r.FindById(ctx, int(insertedId))
}

// ユーザーのデッキ一覧取得
func (r *deckRepository) FindAll(ctx context.Context) ([]*deck.Deck, error) {
	query := db.GetQuery(ctx)

	// ユーザーのデッキ一覧を取得
	deckRows, err := query.FindALl(ctx)
	if err != nil {
		return nil, fmt.Errorf("デッキ一覧取得エラー: %w", err)
	}

	var decks []*deck.Deck
	for _, row := range deckRows {
		// 個別のデッキを取得
		deck, err := r.FindById(ctx, int(row.ID))
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}

	return decks, nil
}

// デッキの詳細取得
func (r *deckRepository) FindById(ctx context.Context, id int) (*deck.Deck, error) {
	query := db.GetQuery(ctx)

	// デッキの基本情報を取得
	deckRow, err := query.FindDeckById(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("デッキが見つかりません")
		}
		return nil, fmt.Errorf("デッキ取得エラー: %w", err)
	}

	// デッキカードを取得
	deckCardRows, err := query.FindDeckCardsByDeckId(ctx, deckRow.ID)
	if err != nil {
		return nil, fmt.Errorf("デッキカード取得エラー: %w", err)
	}

	// メインカード、サブカード、デッキカードをドメインオブジェクトに変換
	var mainCard domain.Card
	var subCard domain.Card
	var deckCards []deck.DeckCard

	// メインカードがある場合
	if deckRow.MainCardID.Valid && deckRow.MainCardTypeID.Valid {
		card, err := r.cardRepository.FindCardById(ctx, int(deckRow.MainCardID.Int64), domain.CardType(deckRow.MainCardTypeID.Int64))
		if err != nil {
			return nil, fmt.Errorf("メインカード取得エラー: %w", err)
		}
		mainCard = card
	}

	// サブカードがある場合
	if deckRow.SubCardID.Valid && deckRow.SubCardTypeID.Valid {
		card, err := r.cardRepository.FindCardById(ctx, int(deckRow.SubCardID.Int64), domain.CardType(deckRow.SubCardTypeID.Int64))
		if err != nil {
			return nil, fmt.Errorf("サブカード取得エラー: %w", err)
		}
		subCard = card
	}

	// デッキカードを変換
	for _, cardRow := range deckCardRows {
		card, err := r.cardRepository.FindCardById(ctx, int(cardRow.CardID), domain.CardType(cardRow.CardTypeID))
		if err != nil {
			return nil, fmt.Errorf("カード取得エラー: %w", err)
		}

		deckCard := deck.NewDeckCard(card, int(cardRow.Quantity))
		deckCards = append(deckCards, *deckCard)
	}

	// デッキオブジェクトを作成（バリデーションをスキップ）
	description := ""
	if deckRow.Description.Valid {
		description = deckRow.Description.String
	}

	// データベースから読み込むときはバリデーションをスキップ
	return deck.NewDeckWithoutValidation(
		int(deckRow.ID),
		deckRow.Name,
		description,
		mainCard,
		subCard,
		deckCards,
	), nil
}

// デッキの更新
func (r *deckRepository) Update(ctx context.Context, d *deck.Deck) error {
	// トランザクション開始
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback() // 明示的にコミットされなければロールバックする

	qtx := dbgen.New(tx)

	// メインカードとサブカードのID、タイプを取得
	var mainCardID sql.NullInt64
	var mainCardTypeID sql.NullInt64
	var subCardID sql.NullInt64
	var subCardTypeID sql.NullInt64

	if d.GetMainCard() != nil {
		mainCardID.Int64 = int64(d.GetMainCard().GetId())
		mainCardID.Valid = true
		mainCardTypeID.Int64 = int64(d.GetMainCard().GetCardType())
		mainCardTypeID.Valid = true
	}

	if d.GetSubCard() != nil {
		subCardID.Int64 = int64(d.GetSubCard().GetId())
		subCardID.Valid = true
		subCardTypeID.Int64 = int64(d.GetSubCard().GetCardType())
		subCardTypeID.Valid = true
	}

	// デッキを更新
	err = qtx.UpdateDeck(ctx, dbgen.UpdateDeckParams{
		ID:             int64(d.GetId()),
		Name:           d.GetName(),
		Description:    sql.NullString{String: d.GetDescription(), Valid: d.GetDescription() != ""},
		MainCardID:     mainCardID,
		MainCardTypeID: mainCardTypeID,
		SubCardID:      subCardID,
		SubCardTypeID:  subCardTypeID,
	})
	if err != nil {
		return fmt.Errorf("デッキ更新エラー: %w", err)
	}

	// 既存のデッキカードをすべて削除
	err = qtx.DeleteDeckCardsByDeckId(ctx, int64(d.GetId()))
	if err != nil {
		return fmt.Errorf("デッキカード削除エラー: %w", err)
	}

	// デッキカードを新たに追加
	for _, card := range d.GetCards() {
		_, err = qtx.CreateDeckCard(ctx, dbgen.CreateDeckCardParams{
			DeckID:     int64(d.GetId()),
			CardID:     int64(card.GetCard().GetId()),
			CardTypeID: int64(card.GetCard().GetCardType()),
			Quantity:   int32(card.GetQuantity()),
		})
		if err != nil {
			return fmt.Errorf("デッキカード作成エラー: %w", err)
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	return nil
}

// デッキの削除
func (r *deckRepository) Delete(ctx context.Context, id int) error {
	query := db.GetQuery(ctx)

	// デッキを削除（カスケード削除によりデッキカードも削除される）
	err := query.DeleteDeck(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("デッキ削除エラー: %w", err)
	}

	return nil
}
