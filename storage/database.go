// package storage

// import (
// 	"habitat-server/models"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func connectToDB() *gorm.DB {
// 	err := godotenv.Load()
// 	if err != nil {
// 		panic("Error loading .env file")
// 	}

// 	dsn := os.Getenv("DB_CONNECTION_STRING")
// 	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if dbError != nil {
// 		log.Panic("error connection to db")
// 	}

// 	DB = db
// 	return db
// }

// func performMigrations(db *gorm.DB) {
// 	db.AutoMigrate(
// 		&models.Conversation{}, // create table containing many side first
// 		&models.Message{},
// 		&models.User{},
// 		&models.Property{},
// 		&models.Review{},
// 		&models.Apartment{},
// 	)
// }

// func InitializeDB() *gorm.DB {
// 	db := connectToDB()
// 	performMigrations(db)
// 	return db
// }

package storage

import (
	"habitat-server/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// connectToDB connects to the database using the connection string from the environment variable.
func connectToDB() *gorm.DB {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	if dsn == "" {
		log.Panic("DB_CONNECTION_STRING environment variable is not set")
	}

	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbError != nil {
		log.Panic("failed to connect to the database:", dbError)
	}

	DB = db
	return db
}

// performMigrations applies the necessary migrations to the database.
func performMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.User{},
		&models.Property{},
		&models.Review{},
		&models.Apartment{},
	)
	if err != nil {
		log.Panic("failed to perform migrations:", err)
	}
}

// InitializeDB initializes the database connection and performs migrations.
func InitializeDB() *gorm.DB {
	db := connectToDB()
	performMigrations(db)
	return db
}
