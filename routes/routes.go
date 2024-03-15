package routes

import (
	"log"
	"net/http"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/controllers"
	"github.com/gorilla/mux"
)

func HandleRequest() {
	r := mux.NewRouter()

	r.HandleFunc("/api/holders", controllers.CreateNewHolder).Methods("Post")     //criar titular
	r.HandleFunc("/api/holders", controllers.GetAllHolders).Methods("Get")        //buscar todos os titulares
	r.HandleFunc("/api/holders/{id}", controllers.ReturnHolder).Methods("Get")    //buscar titular especifico
	r.HandleFunc("/api/holders/{id}", controllers.DeleteHolder).Methods("Delete") //deletar titular

	r.HandleFunc("/api/accounts/{cpf}", controllers.CreateAccount).Methods("Post") //buscar criar conta usando o cpf do titular

	r.HandleFunc("/api/accounts", controllers.GetAllAccounts).Methods("Get") //buscar todas as contas

	r.HandleFunc("/api/accounts/{number}", controllers.CheckAccount).Methods("Get") //buscar conta especifica

	r.HandleFunc("/api/accounts/block/{number}", controllers.BlockAccount).Methods("Patch")   //bloquear conta
	r.HandleFunc("/api/accounts/unlock/{number}", controllers.UnlockAccount).Methods("Patch") //desbloquear conta

	r.HandleFunc("/api/accounts/disable/{number}", controllers.DisableAccount).Methods("Patch") //desativar conta
	r.HandleFunc("/api/accounts/enable/{number}", controllers.EnableAccount).Methods("Patch")   //habilitar conta

	r.HandleFunc("/api/accounts/{number}/statement", controllers.RequestStatement).Methods("GET") //extrato

	r.HandleFunc("/api/accounts/{number}/deposit/{amount}", controllers.DepositAccount).Methods("Patch") // deposito

	r.HandleFunc("/api/accounts/{number}/withdraw/{amount}", controllers.WithdrawAccount).Methods("Patch") // saque

	log.Fatal(http.ListenAndServe(":8000", r))

}
