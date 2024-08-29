package main

import (
	"filend/routes"
	"filend/services"
	"fmt"
)

func main() {

	fmt.Println(services.GenerateOneTimePassword())

	router := routes.SetupRouter()
	router.Static("/ui", "./ui")
	router.Run(":9090")
}
