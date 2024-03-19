package models

import (
	"database/sql"
	"errors"
)

type Account struct {
	Account_id int     `json:"Account_id"`
	Number     string  `json:"number"`
	Fk_holder  string  `json:"holderCpf"`
	Agency     string  `json:"agency"`
	Balance    float64 `json:"balance"`
	Blocked    bool
	Active     bool
}

func InsertAccount(db *sql.DB, account *Account) error {
	query := "INSERT INTO accounts (fk_holder, number, agency, balance, blocked, active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING account_id"
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare SQL statement: " + err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(&account.Fk_holder, &account.Number, &account.Agency, &account.Balance, &account.Blocked, &account.Active).Scan(&account.Account_id)
	if err != nil {
		return errors.New("failed to insert Account into database: " + err.Error())
	}
	return nil
}

func SearchAccount(db *sql.DB, number string) (*Account, error) {
	query := "SELECT * FROM accounts WHERE number = $1"
	var account Account
	err := db.QueryRow(query, number).Scan(&account.Account_id, &account.Fk_holder, &account.Number, &account.Agency, &account.Balance, &account.Blocked, &account.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &account, nil
}

func UpdateAccountBlock(db *sql.DB, number string) (*Account, error) {
	query := "UPDATE accounts SET blocked = true WHERE number = $1"
	var account Account
	err := db.QueryRow(query, number).Scan(&account.Account_id, &account.Fk_holder, &account.Number, &account.Agency, &account.Blocked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &account, nil
}

func UpdateAccountUnlock(db *sql.DB, number string) (*Account, error) {
	query := "UPDATE accounts SET blocked = false WHERE number = $1"
	var account Account
	err := db.QueryRow(query, number).Scan(&account.Account_id, &account.Fk_holder, &account.Number, &account.Agency, &account.Blocked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &account, nil
}

func UpdateAccountActivate(db *sql.DB, number string) (*Account, error) {
	query := "UPDATE accounts SET active = true WHERE number = $1"
	var account Account
	err := db.QueryRow(query, number).Scan(&account.Account_id, &account.Fk_holder, &account.Number, &account.Agency, &account.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &account, nil
}

func UpdateAccountDeactivate(db *sql.DB, number string) (*Account, error) {
	query := "UPDATE accounts SET active = false WHERE number = $1"
	var account Account
	err := db.QueryRow(query, number).Scan(&account.Account_id, &account.Fk_holder, &account.Number, &account.Agency, &account.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &account, nil
}

func UpdateAccountBalance(db *sql.DB, number string, amount float64, operationType string) error {

	currentBalance, err := GetCurrentAccountBalance(db, number)
	if err != nil {
		return err
	}

	var newBalance float64

	if operationType == "deposit" {
		newBalance = currentBalance + amount
	} else {
		newBalance = currentBalance - amount
	}

	query := "UPDATE accounts SET balance = $1 WHERE number = $2"
	_, err = db.Exec(query, newBalance, number)
	if err != nil {
		return errors.New("failed to update account balance in database: " + err.Error())
	}

	return nil
}

func GetCurrentAccountBalance(db *sql.DB, number string) (float64, error) {

	var balance float64
	query := "SELECT balance FROM accounts WHERE number = $1"
	err := db.QueryRow(query, number).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("account not found")
		}
		return 0, errors.New("failed to query account balance from database: " + err.Error())
	}

	return balance, nil
}
