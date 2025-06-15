package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	AppPort   string
	JWTSecret string
	JWTExpiry int
}

var Config AppConfig

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Printf("Warning: Could not convert environment variable %s to int. Using default value %d. Error: %v\n", key, defaultValue, err)
		return defaultValue
	}
	return intValue
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env. Falling back to system environment variables. Error:", err)
	}
	Config = AppConfig{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnvAsInt("DB_PORT", 5432),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", ""),
		DBName:    getEnv("DB_NAME", ""),
		AppPort:   getEnv("APP_PORT", "50001"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpiry: getEnvAsInt("JWT_EXPIRY", 1440),
	}
}
