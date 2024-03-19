package models

import (
	"database/sql"
	"errors"
	"time"
)

type Operation struct {
	Operation_id int       `json:"operation_id"`
	Fk_account   string    `json:"fk_account"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
}

func CreateOperation(db *sql.DB, operation *Operation) error {
	query := "INSERT INTO operations (fk_account, type, amount, status) VALUES ($1, $2, $3, $4) RETURNING operation_id"
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare SQL statement: " + err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(&operation.Fk_account, &operation.Type, &operation.Amount, &operation.Status).Scan(&operation.Operation_id)
	if err != nil {
		return errors.New("failed to insert Account into database: " + err.Error())
	}
	return nil
}

func GetDailyLimit(db *sql.DB, number string, date time.Time) (float64, error) {

	startOfTime := date.Truncate(24 * time.Hour)
	endOfTime := startOfTime.Add(24 * time.Hour)

	var dailyWithdraw float64
	query := "SELECT COALESCE(SUM(amount), 0) FROM operations WHERE fk_account = $1 AND type = 'withdrawal' AND status = 'Successful' AND date >= $2 AND date < $3"
	err := db.QueryRow(query, number, startOfTime, endOfTime).Scan(&dailyWithdraw)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.New("failed to query account balance from database: " + err.Error())
	}

	return dailyWithdraw, nil
}

func GetStatement(db *sql.DB, number string, minDate time.Time, maxDate time.Time) ([]Operation, error) {

	var operations []Operation

	query := "SELECT operation_id, fk_account, type, amount, date, status FROM operations WHERE fk_account = $1 AND date >= $2 AND date < $3"
	rows, err := db.Query(query, number, minDate, maxDate)
	if err != nil {
		return nil, errors.New("failed to query operations from database: " + err.Error())
	}

	defer rows.Close()
	for rows.Next() {

		var operation Operation

		err := rows.Scan(&operation.Operation_id, &operation.Fk_account, &operation.Type, &operation.Date, &operation.Status)
		if err != nil {
			return nil, errors.New("failed to scan operation row: " + err.Error())
		}
		operations = append(operations, operation)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("error while iterating over operation rows: " + err.Error())
	}

	return operations, nil
}
