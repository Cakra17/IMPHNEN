package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Port      string
	DSN       string
	JWTSecret string
}

func Load() Config {
	return Config{
		Port:      getEnvWithDefaultValue("PORT", ":6969"),
		DSN:       getEnvWithDefaultValue("DSN", "postgres://admin:adminsecret@localhost:5432/imphnen?sslmode=disable"),
		JWTSecret: getEnvWithDefaultValue("JWT_SECRET", "secret"),
	}
}

func getEnvWithDefaultValue(key string, defaultValue string) string {
	res := os.Getenv(key)
	if res == "" {
		os.Setenv(key, defaultValue)
		return defaultValue
	}
	return res
}

func ConnectDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	log.Println("Connected to database!!")
	return db
}
