package main

import (
	"fmt"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/routes"
)

func main() {

	fmt.Println("Iniciando o servidor Rest com Go")
	routes.HandleRequest()
}
