package main

import (
	"log"

	"bookstore-project/internal/database"


	"bookstore-project/internal/config"
)

func main() {
	err := config.NewConfig()

	if err != nil {
		panic(err)
	}


	// cnnect to database and verify connection
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")
}
