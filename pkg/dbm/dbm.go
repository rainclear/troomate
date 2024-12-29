package dbm

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rainclear/troomate/pkg/config"
)

var app *config.AppConfig

func NewDbm(a *config.AppConfig) {
	app = a
}

func CloseDb() error {
	return app.Db.Close()
}

func OpenDb() error {
	db, err := sql.Open("sqlite", app.DBPath)
	app.Db = db

	if err != nil {
		return err
	}

	rows, err := db.Query("Select OwnerName, TFSAEligibleDate from Owners;")
	if err != nil {
		return err
	}

	for rows.Next() {
		var ownername string
		var eligibleDate string
		if err = rows.Scan(&ownername, &eligibleDate); err != nil {
			return err
		}
		app.Owner = ownername
		app.Owner += " from "
		app.Owner += eligibleDate
	}

	accounts, err := ListAccounts()
	if err != nil {
		return err
	}
	app.Accounts = accounts

	fi, err := os.Stat(app.DBPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s size: %v\n", app.DBPath, fi.Size())
	return err
}

func ListAccounts() ([]string, error) {
	var accounts []string
	db := app.Db

	rows, err := db.Query("Select id, AccountName, AccountNameAtCRA from Accounts order by AccountName;")
	if err != nil {
		return accounts, err
	}

	for rows.Next() {
		var account_id string
		var accountname string
		var accountname_at_cra string
		if err = rows.Scan(&account_id, &accountname, &accountname_at_cra); err != nil {
			return accounts, err
		}

		account := account_id
		account += ","
		account += accountname
		account += ","
		account += accountname_at_cra

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return accounts, err
	}

	return accounts, err
}
