package models

import (
	"strconv"
)

type Holder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cpf  string `json:"cpf"`
}

var Holders = make(map[int]Holder)

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
