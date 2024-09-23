package main

import (
	"filend/config"
	"filend/routes"
	"filend/services"
	"fmt"
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
	services.DeleteOldFiles(config.DB)

	fmt.Println(services.GenerateOneTimePassword())

	router := routes.SetupRouter()
	router.Static("/ui", "./ui")
	router.Run(":9091")
}
