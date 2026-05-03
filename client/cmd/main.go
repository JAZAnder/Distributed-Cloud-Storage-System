package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/app"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error Loading .env file")
	}

	a := app.App{}
	session.GetSession()

	a.Start()

}
