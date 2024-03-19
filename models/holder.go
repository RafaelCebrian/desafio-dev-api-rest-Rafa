package models

import (
	"database/sql"
	"errors"
	"strconv"
)

type Holder struct {
	holder_id int    `json:"holder_id"`
	Name      string `json:"name"`
	Cpf       string `json:"cpf"`
}

func VerifyCPF(cpf string) bool {

	var n [11]int

	if len(cpf) != 11 {
		return false
	}

	for i := 0; i < 11; i++ {
		n[i], _ = strconv.Atoi(string(cpf[i]))
	}

	for i := 0; i < 10; i++ {
		if n[i] != n[i+1] {
			break
		}
		if i == 9 {
			return false
		}
	}

	var sum int
	for i := 0; i < 9; i++ {
		sum += n[i] * (10 - i)
	}

	sum %= 11
	if sum < 2 {
		n[9] = 0
	} else {
		n[9] = 11 - sum
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += n[i] * (11 - i)
	}

	sum %= 11
	if sum < 2 {
		n[10] = 0
	} else {
		n[10] = 11 - sum
	}

	if n[9] == n[9] && n[10] == n[10] {
		return true
	}

	return false
}

func InsertHolder(db *sql.DB, holder *Holder) error {
	query := "INSERT INTO holders (name, cpf) VALUES ($1, $2) RETURNING holder_id"
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare SQL statement: " + err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(&holder.Name, &holder.Cpf).Scan(&holder.holder_id)
	if err != nil {
		return errors.New("failed to insert holder into database: " + err.Error())
	}
	return nil
}

func DeleteHolder(db *sql.DB, cpf string) error {
	query := "DELETE FROM holders WHERE cpf = $1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare SQL statement: " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(cpf)
	if err != nil {
		return errors.New("failed to delete holder from database: " + err.Error())
	}
	return nil
}

func SearchCPF(db *sql.DB, cpf string) (bool, error) {
	query := "SELECT COUNT(*) FROM holders WHERE cpf = $1"
	var count int
	err := db.QueryRow(query, cpf).Scan(&count)
	if err != nil {
		return false, errors.New("failed to query database: " + err.Error())
	}
	return count > 0, nil
}

func SearchHolder(db *sql.DB, cpf string) (*Holder, error) {
	query := "SELECT * FROM holders WHERE cpf = $1"
	var holder Holder
	err := db.QueryRow(query, cpf).Scan(&holder.holder_id, &holder.Name, &holder.Cpf)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("holder not found")
		}
		return nil, errors.New("failed to query database: " + err.Error())
	}
	return &holder, nil
}
