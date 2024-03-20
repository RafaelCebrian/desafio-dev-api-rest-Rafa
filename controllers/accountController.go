package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/database"
	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/models"
	"github.com/gorilla/mux"
)

func CreateAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	cpf := vars["cpf"]

	if cpf == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Cpf value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database" + err.Error()))
		return
	}

	exists, err := models.SearchCPF(db, cpf)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error checking CPF in database: " + err.Error()))
		return
	}
	if exists == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("This CPF is not registered in the list of holders"))
		return
	}

	accountNumber := strconv.Itoa(rand.Intn(90000000) + 10000000)
	accountAgency := strconv.Itoa(rand.Intn(900) + 100)

	newAccount := models.Account{
		Fk_holder: cpf,
		Number:    accountNumber,
		Agency:    accountAgency,
		Balance:   0.0,
		Blocked:   false,
		Active:    true,
	}

	err = models.InsertAccount(db, &newAccount)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error inserting account into database" + err.Error()))
		return
	}

	rw.WriteHeader(http.StatusCreated)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(newAccount)
	defer db.Close()
}

func ReturnAccount(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	account, err := models.SearchAccount(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(account)
	defer db.Close()
}
func DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("account number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.DeleteAccountDB(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at deleting account from database" + err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Account deleted successfully"))
	defer db.Close()
}
func BlockAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.UpdateAccountBlock(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Account blocked successfully"))
	defer db.Close()

}

func UnlockAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.UpdateAccountUnlock(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Account unlocked successfully"))
	defer db.Close()

}

func DisableAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.UpdateAccountDeactivate(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Account successfully deactivated"))
	defer db.Close()

}

func EnableAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	number := vars["number"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.UpdateAccountActivate(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database" + err.Error()))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Account successfully activated"))
	defer db.Close()
}
