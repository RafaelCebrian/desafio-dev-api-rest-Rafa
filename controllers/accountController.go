package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"math/rand"
	"net/http"

	"github.com/RafaelCebrian/desafio-dev-api-rest-Rafa/models"
	"github.com/gorilla/mux"
)

var lastIDAccount int

func init() {
	lastIDAccount = 0
}

func CreateAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	cpf := vars["cpf"]

	if cpf == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Cpf value is null"))
		return
	}

	// if SearchCPF(cpf) == false {
	// 	rw.WriteHeader(http.StatusNotFound)
	// 	rw.Write([]byte("Cpf not found"))
	// 	return
	// }

	agencyAccount := rand.Intn(10000)

	agencyNumber := rand.Intn(500)

	lastIDAccount++

	newAccount := models.Account{
		ID:        lastIDAccount,
		Number:    agencyNumber,
		Agency:    agencyAccount,
		Balance:   0.0,
		Blocked:   false,
		Active:    true,
		Statement: make([]models.Operation, 0),
	}

	models.Accounts[cpf] = newAccount
	rw.WriteHeader(http.StatusCreated)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(newAccount)

}

func GetAllAccounts(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(models.Accounts)

}

func CheckAccount(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	numberStr := vars["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	var foundAcc bool

	for _, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(account)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}
}

func BlockAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numberStr := vars["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]
			accChange.Blocked = true
			models.Accounts[key] = accChange

			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}

}

func UnlockAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numberStr := vars["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]
			accChange.Blocked = false
			models.Accounts[key] = accChange

			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}

}

func DisableAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numberStr := vars["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]
			accChange.Active = false
			models.Accounts[key] = accChange

			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}

}

func EnableAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numberStr := vars["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]
			accChange.Active = true
			models.Accounts[key] = accChange

			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}
}

func DepositAccount(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	numberStr := vars["number"]
	amountStr := vars["amount"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the amount"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]

			if accChange.Active == false || accChange.Blocked == true {

				rw.WriteHeader(http.StatusForbidden)
				if accChange.Active == false {
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

				currentTime := time.Now()

				formatedTime := currentTime.Format(time.RFC3339)

				operationStatus := "Successful"

				depositOperation := models.Operation{
					Type:   "deposit",
					Amount: amount,
					Date:   formatedTime,
					Status: operationStatus,
				}

				accChange.Statement = append(accChange.Statement, depositOperation)
				accChange.Balance += amount

				models.Accounts[key] = accChange

				rw.WriteHeader(http.StatusOK)
				return
			}

		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}

}
func WithdrawAccount(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	numberStr := vars["number"]
	amountStr := vars["amount"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the amount"))
		return
	}

	var foundAcc bool

	for key, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			accChange := models.Accounts[key]

			if accChange.Active == false || accChange.Blocked == true {
				rw.WriteHeader(http.StatusForbidden)
				return
			} else {
				operationStatus := "Successful"

				if amount <= 0 {
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write([]byte("Insufficient withdrawal amount"))
					return
				}

				currentTime := time.Now()

				formatedTime := currentTime.Format(time.RFC3339)

				verifyTimeDay := currentTime.Truncate(24 * time.Hour)

				dailyLimit, exists := models.DailyLimits[verifyTimeDay]

				if exists == false {
					dailyLimit = models.DailyLimit{
						Date:          verifyTimeDay,
						TotalWithdraw: 0,
					}
				}

				if amount > account.Balance {
					operationStatus = "Failed - Insufficient funds"
				} else if dailyLimit.TotalWithdraw+amount > 2000 {
					operationStatus = "Failed - Daily withdraw limit exceeded"
				} else {
					dailyLimit.TotalWithdraw += amount
					models.DailyLimits[verifyTimeDay] = dailyLimit

					accChange.Balance -= amount
					models.Accounts[key] = accChange
				}

				withdrawOperation := models.Operation{
					Type:   "withdraw",
					Amount: amount,
					Date:   formatedTime,
					Status: operationStatus,
				}

				accChange.Statement = append(accChange.Statement, withdrawOperation)
				models.Accounts[key] = accChange

				rw.WriteHeader(http.StatusOK)
				return
			}
		}
	}
	if !foundAcc {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}
}

func RequestStatement(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numberStr := vars["number"]

	UrlDateMin := req.URL.Query().Get("min")
	UrlDateMax := req.URL.Query().Get("max")

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error in requesting the account address"))
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

	var foundAcc bool

	for _, account := range models.Accounts {
		if account.Number == number {

			foundAcc = true

			var filteredStatement []models.Operation
			for _, operation := range account.Statement {

				operationDate, err := time.Parse(time.RFC3339, operation.Date)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte("Error processing operation date"))
					return
				}

				if operationDate.After(minDate) && operationDate.Before(maxDate) {
					filteredStatement = append(filteredStatement, operation)
				}
			}

			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(filteredStatement)
			return
		}
	}
	if foundAcc == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Account not found"))
	}
}
