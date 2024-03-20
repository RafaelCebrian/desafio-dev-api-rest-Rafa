package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/database"
	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/models"
	"github.com/gorilla/mux"
)

func DepositAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	number := vars["number"]
	amountStr := vars["amount"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the amount"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database" + err.Error()))
		return
	}

	account, err := models.SearchAccount(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting account from database"))
		return
	}

	if account.Active == false {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Inactive account"))
		return
	}

	if account.Blocked == true {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Blocked account"))
		return
	}

	if amount < 0 {

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Negative deposit amount"))
		return
	}

	operationStatus := "Successful"

	newOperation := models.Operation{
		Fk_account: number,
		Type:       "deposit",
		Amount:     amount,
		Status:     operationStatus,
	}

	err = models.CreateOperation(db, &newOperation)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at creating an deposit operation in database"))
		return
	}

	err = models.UpdateAccountBalance(db, number, amount, newOperation.Type)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at updating the account balance in database"))
		return
	}
	defer db.Close()
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Successfully deposit: " + amountStr))
}

func WithdrawAccount(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	number := vars["number"]
	amountStr := vars["amount"]

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the amount"))
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database" + err.Error()))
		return
	}

	account, err := models.SearchAccount(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at getting account from database"))
		return
	}

	if account.Active == false {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Inactive account"))
		return
	}

	if account.Blocked == true {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Blocked account"))
		return
	}

	if amount < 0 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Negative withdrawal amount"))
		return
	}

	totalWithdraws, err := models.GetDailyLimit(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to get total withdrawals for the day" + err.Error()))
		return
	}

	dailyLimit := 2000.0
	if totalWithdraws+amount > dailyLimit {
		newOperation := models.Operation{
			Fk_account: number,
			Type:       "withdrawal",
			Amount:     amount,
			Status:     "Failed - Daily withdrawal limit exceeded",
		}

		err := models.CreateOperation(db, &newOperation)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error at creating a withdrawal operation in database"))
			return
		}

		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Daily withdrawal limit exceeded"))
		return
	}

	currentBalance, err := models.GetCurrentAccountBalance(db, number)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at geting the current balance"))
		return
	}

	if currentBalance-amount < 0 {
		newOperation := models.Operation{
			Fk_account: number,
			Type:       "withdrawal",
			Amount:     amount,
			Status:     "Failed - insufficient funds",
		}

		err := models.CreateOperation(db, &newOperation)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error at creating a withdrawal operation in database"))
			return
		}

		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Insufficient funds"))
		return
	}

	newOperation := models.Operation{
		Fk_account: number,
		Type:       "withdrawal",
		Amount:     amount,
		Status:     "Successful",
	}

	err = models.CreateOperation(db, &newOperation)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at creating a withdrawal operation in database"))
		return
	}

	err = models.UpdateAccountBalance(db, number, amount, newOperation.Type)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at updating the account balance in database"))
		return
	}

	defer db.Close()
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Successfully withdrawn: " + amountStr))
}
func RequestStatement(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	number := vars["number"]

	UrlDateMin := req.URL.Query().Get("min")
	UrlDateMax := req.URL.Query().Get("max")

	if number == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("number value is null"))
		return
	}

	minDate, err := time.Parse(time.RFC3339, UrlDateMin)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error getting the minimum date"))
		return
	}

	maxDate, err := time.Parse(time.RFC3339, UrlDateMax)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error getting maximum date"))
		return

	}

	db, err := database.ConnectDB()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to connect to database" + err.Error()))
		return
	}

	operations, err := models.GetStatement(db, number, minDate, maxDate)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error at retrieving account statement from database" + err.Error()))
		return
	}

	defer db.Close()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(operations)
	rw.WriteHeader(http.StatusOK)

}
