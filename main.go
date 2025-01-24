package main

import (
	"filend/config"
	"filend/routes"
	"filend/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.ConnectDatabase()
	config.ConnectMinio()

	go services.StartScheduler(config.DB)

	router := routes.SetupRouter()
	router.Static("/ui", "./ui")
	router.Run(":9091")
}
