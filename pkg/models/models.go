package models

type Account struct {
	ID               int64
	AccountName      string
	Institution      string
	AccountNumber    string
	AccountNameAtCRA string
	AccountType      string
	AccountPurpose   string
	Notes            string
}
