package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBService interface defines the Connect method
type DBService interface {
	Connect() (*gorm.DB, error)
}

// PostGresDB  implements the DBService interface
type PostgresDB struct{}

// Service for instance
func NewDBService() DBService {
	return &PostgresDB{}
}

// Connect opens a connection to the PostGreSQL database
func (p *PostgresDB) Connect() (*gorm.DB, error) {
	// Construct the Data Source Name (DSN) string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Config.DBHost,
		Config.DBUser,
		Config.DBPassword,
		Config.DBName,
		Config.DBPort,
	)

	// Log the DSN (for debugging purposes only, remove in production)
	log.Println("Database DSN:", dsn)

	// Try connecting to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("❌ Failed to connect to PostgreSQL:", err)
		return nil, err
	}

	// Log successful connection
	log.Println("✅ Successfully connected to PostgreSQL database.")
	return db, nil
}
