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

// 	dsn := os.Getenv("DATABASE_URL")
// 	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if dbError != nil {
// 		log.Panic("error connection to db")
// 	}
// 	if err != nil {
//         log.Printf("[error] failed to initialize database, got error %v", err)
//         log.Panic("error connection to db")
//     }

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

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDB() *gorm.DB {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	// Get the database URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Panic("DATABASE_URL is not set in the environment variables")
	}

	// Connect to the database
	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbError != nil {
		log.Printf("[error] failed to initialize database, got error %v", dbError)
		log.Panic("Error connecting to the database")
	}

	// Assign the db to the global variable
	DB = db
	return db
}

func performMigrations(db *gorm.DB) {
	// Perform database migrations
	db.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.User{},
		&models.Property{},
		&models.Review{},
		&models.Apartment{},
	)
}

func InitializeDB() *gorm.DB {
	db := connectToDB()
	performMigrations(db)
	return db
}


