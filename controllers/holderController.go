package controllers

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/models"
	"github.com/gorilla/mux"
)

var lastIDHolder int

func init() {
	lastIDHolder = 0
}

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
	if SearchCPF(newHolder.Cpf) == true {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("CPF already registered"))
		return
	}

	lastIDHolder++
	newHolder.ID = lastIDHolder
	models.Holders[lastIDHolder] = newHolder

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(newHolder)
}

func DeleteHolder(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in ID value"))
		return
	}
	_, found := models.Holders[id]
	if found == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Holder not found"))
		return
	}

	delete(models.Holders, id)
	rw.WriteHeader(http.StatusOK)
}

func GetAllHolders(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(models.Holders)

}

func ReturnHolder(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in ID value"))
		return
	}

	holder, found := models.Holders[id]

	if found == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Holder not found"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(holder)

}

func SearchCPF(cpf string) bool {
	for _, holder := range models.Holders {
		if holder.Cpf == cpf {
			return true
		}
	}
	return false
}
