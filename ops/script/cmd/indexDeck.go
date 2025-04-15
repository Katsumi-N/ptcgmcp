/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// インデックスするデッキのデータ構造
type DeckDocument struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	MainCard    *CardInfo  `json:"main_card,omitempty"`
	SubCard     *CardInfo  `json:"sub_card,omitempty"`
	Cards       []DeckCard `json:"cards"`
}

// カード情報の構造体
type CardInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
}

// デッキカードの構造体
type DeckCard struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
	Quantity int    `json:"quantity"`
}

// indexDeckCmd represents the indexDeck command
var indexDeckCmd = &cobra.Command{
	Use:   "index-deck",
	Short: "Index decks to Meilisearch",
	Long:  `Index decks with their cards to Meilisearch for search functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Indexing decks to Meilisearch...")

		// 設定を読み込む
		mysqlConfig := MySQLConfig{
			User:     viper.GetString("mysql.user"),
			Password: viper.GetString("mysql.password"),
			Host:     viper.GetString("mysql.host"),
			Port:     viper.GetString("mysql.port"),
			DBName:   viper.GetString("mysql.dbname"),
		}

		meiliConfig := MeilisearchConfig{
			Host: viper.GetString("meilisearch.host"),
			Key:  viper.GetString("meilisearch.key"),
		}

		// デフォルト値を設定
		if mysqlConfig.Host == "" {
			mysqlConfig.Host = "localhost"
		}
		if mysqlConfig.Port == "" {
			mysqlConfig.Port = "3306"
		}
		if mysqlConfig.User == "" {
			mysqlConfig.User = "root"
		}
		if mysqlConfig.DBName == "" {
			mysqlConfig.DBName = "ptcgmcpdb"
		}

		if meiliConfig.Host == "" {
			meiliConfig.Host = "http://localhost:7700"
		}
		if meiliConfig.Key == "" {
			meiliConfig.Key = "DevelopmentMasterKey"
		}

		// データベースに接続
		db, err := connectDB(mysqlConfig)
		if err != nil {
			log.Fatalf("データベース接続エラー: %v", err)
		}
		defer db.Close()

		IndexDeck(db, meiliConfig)
	},
}

func init() {
	rootCmd.AddCommand(indexDeckCmd)

	// コマンドフラグの定義
	indexDeckCmd.Flags().String("mysql-host", "localhost", "MySQL host")
	indexDeckCmd.Flags().String("mysql-port", "3306", "MySQL port")
	indexDeckCmd.Flags().String("mysql-user", "root", "MySQL user")
	indexDeckCmd.Flags().String("mysql-password", "pass", "MySQL password")
	indexDeckCmd.Flags().String("mysql-dbname", "ptcgmcpdb", "MySQL database name")

	indexDeckCmd.Flags().String("meilisearch-host", "http://localhost:7700", "Meilisearch host")
	indexDeckCmd.Flags().String("meilisearch-key", "DevelopmentMasterKey", "Meilisearch API key")

	// viperとフラグをバインド
	viper.BindPFlag("mysql.host", indexDeckCmd.Flags().Lookup("mysql-host"))
	viper.BindPFlag("mysql.port", indexDeckCmd.Flags().Lookup("mysql-port"))
	viper.BindPFlag("mysql.user", indexDeckCmd.Flags().Lookup("mysql-user"))
	viper.BindPFlag("mysql.password", indexDeckCmd.Flags().Lookup("mysql-password"))
	viper.BindPFlag("mysql.dbname", indexDeckCmd.Flags().Lookup("mysql-dbname"))

	viper.BindPFlag("meilisearch.host", indexDeckCmd.Flags().Lookup("meilisearch-host"))
	viper.BindPFlag("meilisearch.key", indexDeckCmd.Flags().Lookup("meilisearch-key"))
}

func IndexDeck(db *sql.DB, meiliConfig MeilisearchConfig) {
	client := meilisearch.New(meiliConfig.Host, meilisearch.WithAPIKey(meiliConfig.Key))
	fmt.Println("Connected to Meilisearch at", meiliConfig.Host)

	index := client.Index("decks")

	// すべてのデッキを取得
	rows, err := db.Query(`
		SELECT 
			d.id, 
			d.name, 
			d.description,
			d.main_card_id, 
			d.main_card_type_id,
			d.sub_card_id, 
			d.sub_card_type_id
		FROM 
			decks d
	`)
	if err != nil {
		log.Fatalf("デッキデータ取得エラー: %v", err)
	}
	defer rows.Close()

	var decks []DeckDocument
	for rows.Next() {
		var deck DeckDocument
		var description sql.NullString
		var mainCardID, mainCardTypeID, subCardID, subCardTypeID sql.NullInt64

		// デッキの基本情報を取得
		err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&description,
			&mainCardID,
			&mainCardTypeID,
			&subCardID,
			&subCardTypeID,
		)
		if err != nil {
			log.Printf("デッキデータ読み取りエラー: %v", err)
			continue
		}

		if description.Valid {
			deck.Description = description.String
		}

		// メインカードの情報を取得
		if mainCardID.Valid && mainCardTypeID.Valid {
			mainCard, err := getCardInfo(db, mainCardID.Int64, mainCardTypeID.Int64)
			if err != nil {
				log.Printf("メインカード情報取得エラー: %v", err)
			} else {
				deck.MainCard = mainCard
			}
		}

		// サブカードの情報を取得
		if subCardID.Valid && subCardTypeID.Valid {
			subCard, err := getCardInfo(db, subCardID.Int64, subCardTypeID.Int64)
			if err != nil {
				log.Printf("サブカード情報取得エラー: %v", err)
			} else {
				deck.SubCard = subCard
			}
		}

		// デッキに含まれるカードを取得
		deckCards, err := getDeckCards(db, deck.ID)
		if err != nil {
			log.Printf("デッキカード情報取得エラー: %v", err)
		}
		deck.Cards = deckCards

		decks = append(decks, deck)
	}

	if len(decks) == 0 {
		fmt.Println("インデックスするデッキがありません")
		return
	}

	// Meilisearchにデッキをインデックス
	_, err = index.AddDocuments(decks)
	if err != nil {
		log.Fatalf("デッキインデックス作成エラー: %v", err)
	}

	// 検索可能なフィールドを設定
	searchableAttributes := []string{"name", "description", "main_card.name", "sub_card.name", "cards.name"}
	_, err = index.UpdateSearchableAttributes(&searchableAttributes)
	if err != nil {
		log.Printf("検索可能フィールド設定エラー: %v", err)
	}

	// ソート可能なフィールドを設定
	sortableAttributes := []string{"id", "name"}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		log.Printf("ソート可能フィールド設定エラー: %v", err)
	}

	fmt.Printf("合計 %d 件のデッキをインデックスしました\n", len(decks))
}

// カード情報を取得する関数
func getCardInfo(db *sql.DB, cardID int64, cardTypeID int64) (*CardInfo, error) {
	var category string
	var tableName string

	// カードタイプ(1:pokemon, 2:trainer, 3:energy)に基づいてテーブルとカテゴリを決定
	switch cardTypeID {
	case 1:
		tableName = "pokemons"
		category = "pokemon"
	case 2:
		tableName = "trainers"
		category = "trainer"
	case 3:
		tableName = "energies"
		category = "energy"
	default:
		return nil, fmt.Errorf("未知のカードタイプ: %d", cardTypeID)
	}

	// カード情報を取得するSQLクエリを作成
	query := fmt.Sprintf("SELECT id, name, image_url FROM %s WHERE id = ?", tableName)

	var card CardInfo
	err := db.QueryRow(query, cardID).Scan(&card.ID, &card.Name, &card.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("カード情報取得エラー: %w", err)
	}

	card.Category = category
	return &card, nil
}

// デッキに含まれるカードを取得する関数
func getDeckCards(db *sql.DB, deckID int64) ([]DeckCard, error) {
	rows, err := db.Query(`
		SELECT
			dc.card_id,
			dc.card_type_id,
			dc.quantity
		FROM
			deck_cards dc
		WHERE
			dc.deck_id = ?
	`, deckID)
	if err != nil {
		return nil, fmt.Errorf("デッキカード取得エラー: %w", err)
	}
	defer rows.Close()

	var cards []DeckCard
	for rows.Next() {
		var cardID int64
		var cardTypeID int64
		var quantity int

		err := rows.Scan(&cardID, &cardTypeID, &quantity)
		if err != nil {
			log.Printf("デッキカード読み取りエラー: %v", err)
			continue
		}

		// カード情報を取得
		cardInfo, err := getCardInfo(db, cardID, cardTypeID)
		if err != nil {
			log.Printf("カード情報取得エラー: %v", err)
			continue
		}

		deckCard := DeckCard{
			ID:       cardInfo.ID,
			Name:     cardInfo.Name,
			Category: cardInfo.Category,
			ImageURL: cardInfo.ImageURL,
			Quantity: quantity,
		}

		cards = append(cards, deckCard)
	}

	return cards, nil
}
