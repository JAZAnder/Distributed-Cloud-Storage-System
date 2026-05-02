package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	. "github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/app"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"

)

// Entry Point
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error Loading .env file")
		
	}

	a := App{}

	database.GetDatabase()

	a.Initialize()

	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	a.Run(port)
}
