package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func Connection() *gorm.DB {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connection := "host=" + dbhost + " port=" + dbport + " user=" + dbuser +
		" dbname=" + dbname + " password=" + dbpassword + " sslmode=disable"

	log.Println(connection)

	// Open the connection to the database
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return db
}
