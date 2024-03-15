package models

import "time"

type Account struct {
	ID        int     `json:"id"`
	Number    int     `json:"number"`
	Agency    int     `json:"agency"`
	Balance   float64 `json:"balance"`
	Blocked   bool
	Active    bool
	Statement []Operation `json:"statement"`
}

type Operation struct {
	// ID     int     `json:"id"` ?
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	Status string  `json:"status"`
}

type DailyLimit struct {
	Date          time.Time
	TotalWithdraw float64
}

var DailyLimits = make(map[time.Time]DailyLimit)

var Accounts = make(map[string]Account)
