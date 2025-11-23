package utils

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database driver => ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect database => ", err)
	}
	log.Println("Connected to database")
	return db
}
