package database

import (
	"FM_techincaltest/app"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DBClient struct {
	DB *sql.DB
}

func InitDB() (*DBClient, error) {
	app.LoadConfig()
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		app.Config.DBHost,
		app.Config.DBPort,
		app.Config.DBUser,
		app.Config.DBPass,
		app.Config.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to PostgreSQL database on port 5432")

	return &DBClient{DB: db}, nil
}

func (client *DBClient) Close() {
	if client != nil && client.DB != nil {
		err := client.DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v\n", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
