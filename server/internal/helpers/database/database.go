package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/user"

)

var once sync.Once

var (
	db *gorm.DB
)

// Singleton Database Connection
func GetDatabase() *gorm.DB {

	once.Do(func() {
		connect()
		migrateTables()
		seedData()
	})

	return db
}

func connect() {

	user := "DB_USERNAME"
	password := "DB_PASSWORD"
	dbServer := "DB_SERVER"
	dbName := "DB_NAME"
	dbPort := "DB_PORT"

	err := godotenv.Load()
	if err != nil {
		log.Println("Error Loading .env file in database")

	} else {
		user = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		dbServer = os.Getenv("DB_SERVER")
		dbName = os.Getenv("DB_NAME")
		dbPort = os.Getenv("DB_PORT")
	}
	if user == "" {log.Println("Warning: env:DEFAULT_USERNAME NOT SET")}
	if password == "" {log.Println("Warning: env:DB_PASSWORD NOT SET")}
	if dbServer == "" {log.Println("Warning: env:DB_SERVER NOT SET")}
	if dbName == "" {log.Println("Warning: env:DB_NAME NOT SET")}
	if dbPort == "" {log.Println("Warning: env:DB_PORT NOT SET")}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbServer, user, password, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	} else {
		log.Println("Connected to Database")
	}
}

func migrateTables() {

	fmt.Println("Migrating database schema...")
	err := db.AutoMigrate(&user.User{}, &securityLog.SecurityLog{})
	if err != nil {
		log.Panicln("failed to Migrating to database: %v", err)
	} else {
		log.Println("Migrated to Database")
	}
}

func seedData() {

	var count int64
	db.Model(&user.User{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding initial administrative user...")
		logEntry := securityLog.SecurityLog{
			Principal: "System",
			Action:    "CreateUser",
		}

		seedUser := user.User{}

		seedUser.Username = os.Getenv("DEFAULT_USERNAME")
		if seedUser.Username == "" {
			seedUser.Username = "Admin"
			log.Println("env:DEFAULT_USERNAME not set defaulting to Admin")
		}

		ptPassword := os.Getenv("DEFAULT_PASSWORD")
		if ptPassword == "" {
			ptPassword = "password"
			log.Println("env:DEFAULT_PASSWORD not set defaulting to password")
		}
		password, err := bcrypt.GenerateFromPassword([]byte(ptPassword), 12)
		if err != nil {
			logEntry.Details = "Failed:" + err.Error()
			db.Create(&logEntry)
			fmt.Println(err)
		}
		seedUser.PasswordHash = string(password)

		result := db.Create(&seedUser)
		logEntry.ResourceID = strconv.FormatUint(uint64(seedUser.ID), 10)

		if result.Error != nil {
			logEntry.Details = "Failed:" + err.Error()
			db.Create(&logEntry)
			log.Fatalf("could not seed user: %v", result.Error)
		}
		db.Create(&logEntry)

		fmt.Printf("User '%s' seeded successfully with ID %d.\n", seedUser.Username, seedUser.ID)
	} else {
		fmt.Println("User records already exist; skipping seed process.")
	}
}
