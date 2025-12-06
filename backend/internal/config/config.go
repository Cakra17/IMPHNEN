package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Port                string
	DSN                 string
	JWTSecret           string
	KolosalApiKey       string
	CloudinaryName      string
	CloudinaryApiKey    string
	CLoudinaryApiSecret string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	log.Println("Env setup finished")
	return Config{
		Port:                os.Getenv("PORT"),
		DSN:                 os.Getenv("DSN"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		KolosalApiKey:       os.Getenv("KOLOSAL_API_KEY"),
		CloudinaryName:      os.Getenv("CLOUDINARY_NAME"),
		CloudinaryApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CLoudinaryApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
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
