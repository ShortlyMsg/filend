package main

import (
	"filend/routes"
	"filend/services"
	"fmt"
)

func main() {

	fmt.Println(services.GenerateOneTimePassword())

	router := routes.SetupRouter()
	router.Run(":8080")
}
