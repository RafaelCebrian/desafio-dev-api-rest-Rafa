package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/database"
	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/models"
	"github.com/gorilla/mux"
)

func CreateNewHolder(rw http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var newHolder models.Holder

	err := decoder.Decode(&newHolder)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in request body"))
		return
	}

	if models.VerifyCPF(newHolder.Cpf) == false {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid CPF"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database" + err.Error()))
		return
	}

	cpf := newHolder.Cpf
	exists, err := models.SearchCPF(db, cpf)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error checking CPF in database: " + err.Error()))
		return
	}
	if exists == true {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("CPF already registered"))
		return
	}

	err = models.InsertHolder(db, &newHolder)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error inserting holder into database" + err.Error()))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(newHolder)

	defer db.Close()
}

func DeleteHolder(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cpf := vars["cpf"]

	if cpf == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cpf value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	err = models.DeleteHolder(db, cpf)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at deleting holder from database"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	defer db.Close()
}

func ReturnHolder(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	cpf := vars["cpf"]

	if cpf == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cpf value is null"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database"))
		return
	}

	holder, err := models.SearchHolder(db, cpf)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting holder from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(holder)
	defer db.Close()

}
