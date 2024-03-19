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

	if account.Active == false || account.Blocked == true {

		rw.WriteHeader(http.StatusForbidden)
		if account.Active == false {
			rw.Write([]byte("Inactive account"))
		} else {
			rw.Write([]byte("Blocked account"))
		}
		return
	} else {

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

		err := models.CreateOperation(db, &newOperation)
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
		rw.WriteHeader(http.StatusOK)
	}
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
		rw.Write([]byte("Error at geting account from database"))
		return
	}
	if account.Active == false || account.Blocked == true {

		rw.WriteHeader(http.StatusForbidden)
		if account.Active == false {
			rw.Write([]byte("Inactive account"))
		} else {
			rw.Write([]byte("Blocked account"))
		}
		return
	} else {

		var operationStatus string

		newOperation := models.Operation{
			Fk_account: number,
			Type:       "withdraw",
			Amount:     amount,
			Status:     operationStatus,
		}

		if amount < 0 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Negative withdraw amount"))
			return
		}

		totalWithdraws, err := models.GetDailyLimit(db, number, time.Now())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Failed to get total withdrawals for the day" + err.Error()))
			return
		}

		dailyLimit := 2000.0
		if totalWithdraws+amount > dailyLimit {

			err := models.CreateOperation(db, &newOperation)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Error at creating an withdraw operation in database"))
				return
			}

			newOperation.Status = "Failed - daily withdrawal limit exceeded"

			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte("Daily withdrawal limit exceeded"))

		} else {

			err := models.CreateOperation(db, &newOperation)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Error at creating an withdraw operation in database"))
				return
			}

			err = models.UpdateAccountBalance(db, number, amount, newOperation.Type)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Error at updating the account balance in database"))
				return
			}

			rw.WriteHeader(http.StatusOK)
		}

	}
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
		rw.Write([]byte("Error at retrieving account statement from database"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(operations)
	rw.WriteHeader(http.StatusOK)

}
