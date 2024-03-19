package main

import (
	"fmt"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/database"
	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/routes"
)

func main() {

	fmt.Println("Starting Rest server with Go")
	database.ConnectDB()
	routes.HandleRequest()
}
