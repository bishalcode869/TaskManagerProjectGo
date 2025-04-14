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
	dsn := fmt.Sprintf(

		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Config.DBHost,
		Config.DBUser,
		Config.DBPassword,
		Config.DBName,
		Config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Successfully connected to PostgreSQL database.")
	return db, nil
}
