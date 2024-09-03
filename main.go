package main

import (
	"filend/config"
	"filend/routes"
	"filend/services"
	"fmt"
)

func main() {

	config.ConnectDatabase()

	fmt.Println(services.GenerateOneTimePassword())

	router := routes.SetupRouter()
	router.Static("/ui", "./ui")
	router.Run(":9090")
}
