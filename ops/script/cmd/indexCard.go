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

// MySQLConfig holds the database connection settings
type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// MeilisearchConfig holds the Meilisearch settings
type MeilisearchConfig struct {
	Host string
	Key  string
}

type Pokemon struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	EnergyType         string   `json:"energy_type"`
	ImageURL           string   `json:"image_url"`
	HP                 int64    `json:"hp"`
	Ability            string   `json:"ability,omitempty"`
	AbilityDescription string   `json:"ability_description,omitempty"`
	Regulation         string   `json:"regulation"`
	Expansion          string   `json:"expansion"`
	Attacks            []Attack `json:"attacks,omitempty"`
}

type Attack struct {
	Name           string `json:"name"`
	RequiredEnergy string `json:"required_energy"`
	Damage         string `json:"damage,omitempty"`
	Description    string `json:"description,omitempty"`
}

type Trainer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	TrainerType string `json:"trainer_type"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}

type Energy struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Regulation  string `json:"regulation"`
	Expansion   string `json:"expansion"`
}

// indexCardCmd represents the indexCard command
var indexCardCmd = &cobra.Command{
	Use:   "index-card",
	Short: "index card to Meilisearch",
	Long:  `index pokemons/trainers/energies to Meilisearch.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Indexing cards to Meilisearch...")

		// Load configuration
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

		// Default values if not found in configuration
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

		// Connect to database
		db, err := connectDB(mysqlConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		IndexCard(db, meiliConfig)
	},
}

func init() {
	rootCmd.AddCommand(indexCardCmd)

	// Define flags for the command
	indexCardCmd.Flags().String("mysql-host", "localhost", "MySQL host")
	indexCardCmd.Flags().String("mysql-port", "3306", "MySQL port")
	indexCardCmd.Flags().String("mysql-user", "root", "MySQL user")
	indexCardCmd.Flags().String("mysql-password", "pass", "MySQL password")
	indexCardCmd.Flags().String("mysql-dbname", "ptcgmcpdb", "MySQL database name")

	indexCardCmd.Flags().String("meilisearch-host", "http://localhost:7700", "Meilisearch host")
	indexCardCmd.Flags().String("meilisearch-key", "DevelopmentMasterKey", "Meilisearch API key")

	// Bind flags with viper
	viper.BindPFlag("mysql.host", indexCardCmd.Flags().Lookup("mysql-host"))
	viper.BindPFlag("mysql.port", indexCardCmd.Flags().Lookup("mysql-port"))
	viper.BindPFlag("mysql.user", indexCardCmd.Flags().Lookup("mysql-user"))
	viper.BindPFlag("mysql.password", indexCardCmd.Flags().Lookup("mysql-password"))
	viper.BindPFlag("mysql.dbname", indexCardCmd.Flags().Lookup("mysql-dbname"))

	viper.BindPFlag("meilisearch.host", indexCardCmd.Flags().Lookup("meilisearch-host"))
	viper.BindPFlag("meilisearch.key", indexCardCmd.Flags().Lookup("meilisearch-key"))
}

func connectDB(config MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func IndexCard(db *sql.DB, config MeilisearchConfig) {
	client := meilisearch.New(config.Host, meilisearch.WithAPIKey(config.Key))

	fmt.Println("Connected to Meilisearch at", config.Host)

	IndexPokemon(client, db)
	IndexTrainer(client, db)
	IndexEnergy(client, db)

	fmt.Println("Indexing complete!")
}

func IndexPokemon(client meilisearch.ServiceManager, db *sql.DB) {
	fmt.Println("Indexing Pokémon cards...")
	index := client.Index("pokemons")

	// Get all Pokémon cards from database
	rows, err := db.Query(`SELECT id, name, energy_type, image_url, hp, 
		ability, ability_description, regulation, expansion 
		FROM pokemons`)
	if err != nil {
		log.Fatalf("Failed to query Pokémon data: %v", err)
	}
	defer rows.Close()

	var pokemons []Pokemon
	for rows.Next() {
		var p Pokemon
		var ability, abilityDesc sql.NullString

		err := rows.Scan(&p.ID, &p.Name, &p.EnergyType, &p.ImageURL, &p.HP,
			&ability, &abilityDesc, &p.Regulation, &p.Expansion)
		if err != nil {
			log.Printf("Error scanning Pokémon row: %v", err)
			continue
		}

		if ability.Valid {
			p.Ability = ability.String
		}
		if abilityDesc.Valid {
			p.AbilityDescription = abilityDesc.String
		}

		// Get attacks for this Pokémon
		p.Attacks = getPokemonAttacks(db, p.ID)
		pokemons = append(pokemons, p)
	}

	if len(pokemons) == 0 {
		fmt.Println("No Pokémon found in database")
		return
	}

	// Index Pokémon cards
	_, err = index.AddDocuments(pokemons)
	if err != nil {
		log.Fatalf("Failed to index Pokémon data: %v", err)
	}

	sortableAttributes := []string{"id"}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		log.Fatalf("Failed to update sortable attributes: %v", err)
	}

	fmt.Printf("Successfully indexed %d Pokémon cards\n", len(pokemons))
}

func getPokemonAttacks(db *sql.DB, pokemonID int64) []Attack {
	rows, err := db.Query(`SELECT name, required_energy, damage, description 
		FROM pokemon_attacks 
		WHERE pokemon_id = ?`, pokemonID)
	if err != nil {
		log.Printf("Error querying attacks for Pokémon ID %d: %v", pokemonID, err)
		return nil
	}
	defer rows.Close()

	var attacks []Attack
	for rows.Next() {
		var a Attack
		var damage, desc sql.NullString

		err := rows.Scan(&a.Name, &a.RequiredEnergy, &damage, &desc)
		if err != nil {
			log.Printf("Error scanning attack row: %v", err)
			continue
		}

		if damage.Valid {
			a.Damage = damage.String
		}
		if desc.Valid {
			a.Description = desc.String
		}

		attacks = append(attacks, a)
	}

	return attacks
}

func IndexTrainer(client meilisearch.ServiceManager, db *sql.DB) {
	fmt.Println("Indexing Trainer cards...")
	index := client.Index("trainers")

	// Get all Trainer cards from database
	rows, err := db.Query(`SELECT id, name, trainer_type, image_url, description, regulation, expansion 
		FROM trainers`)
	if err != nil {
		log.Fatalf("Failed to query Trainer data: %v", err)
	}
	defer rows.Close()

	var trainers []Trainer
	for rows.Next() {
		var t Trainer
		err := rows.Scan(&t.ID, &t.Name, &t.TrainerType, &t.ImageURL, &t.Description, &t.Regulation, &t.Expansion)
		if err != nil {
			log.Printf("Error scanning Trainer row: %v", err)
			continue
		}

		trainers = append(trainers, t)
	}

	if len(trainers) == 0 {
		fmt.Println("No Trainers found in database")
		return
	}

	// Index Trainer cards
	_, err = index.AddDocuments(trainers)
	if err != nil {
		log.Fatalf("Failed to index Trainer data: %v", err)
	}

	sortableAttributes := []string{"id"}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		log.Fatalf("Failed to update sortable attributes: %v", err)
	}

	fmt.Printf("Successfully indexed %d Trainer cards\n", len(trainers))
}

func IndexEnergy(client meilisearch.ServiceManager, db *sql.DB) {
	fmt.Println("Indexing Energy cards...")
	index := client.Index("energies")

	// Get all Energy cards from database
	rows, err := db.Query(`SELECT id, name, image_url, description, regulation, expansion 
		FROM energies`)
	if err != nil {
		log.Fatalf("Failed to query Energy data: %v", err)
	}
	defer rows.Close()

	var energies []Energy
	for rows.Next() {
		var e Energy
		err := rows.Scan(&e.ID, &e.Name, &e.ImageURL, &e.Description, &e.Regulation, &e.Expansion)
		if err != nil {
			log.Printf("Error scanning Energy row: %v", err)
			continue
		}

		energies = append(energies, e)
	}

	if len(energies) == 0 {
		fmt.Println("No Energy cards found in database")
		return
	}

	// Index Energy cards
	_, err = index.AddDocuments(energies)
	if err != nil {
		log.Fatalf("Failed to index Energy data: %v", err)
	}

	sortableAttributes := []string{"id"}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		log.Fatalf("Failed to update sortable attributes: %v", err)
	}

	fmt.Printf("Successfully indexed %d Energy cards\n", len(energies))
}
