package dbm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
	defer rows.Close()

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

	if err := rows.Err(); err != nil {
		log.Fatal(err)
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

	rows, err := db.Query("Select id, OwnerId, AccountName, AccountNameAtCRA from Accounts order by AccountName;")
	if err != nil {
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var account_id string
		var owner_id string
		var accountname string
		var accountname_at_cra string
		if err = rows.Scan(&account_id, &owner_id, &accountname, &accountname_at_cra); err != nil {
			return accounts, err
		}

		account := account_id
		account += ","
		account += owner_id
		account += ","
		account += accountname
		account += ","
		account += accountname_at_cra

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return accounts, err
}

func AddNewAccount(accountname string, accountname_at_cra string) error {
	db := app.Db

	rows, err := db.Query("Select id from Accounts order by id Desc Limit 1;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var next_account_id = 0
	for rows.Next() {
		var account_id string
		if err = rows.Scan(&account_id); err != nil {
			return err
		}

		max_account_id, _ := strconv.Atoi(account_id)
		next_account_id = 1 + max_account_id
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if next_account_id > 0 {
		fmt.Printf("Insert Into Accounts (id, OwnerId, AccountName, AccountNameAtCRA) Values(%d,%d,'%s','%s');\n", next_account_id, 1, accountname, accountname_at_cra)
		stmt, err := db.Prepare("Insert Into Accounts (id, OwnerId, AccountName, AccountNameAtCRA) Values(?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(next_account_id, 1, accountname, accountname_at_cra)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = sql.ErrNoRows
	}

	return err
}

func UpdateAnAccount(accountid_ int, accountname_ string, accountname_at_cra_ string) error {
	db := app.Db

	rows, err := db.Query("Select AccountName, AccountNameAtCRA from Accounts where id=?;", accountid_)
	if err != nil {
		return err
	}
	defer rows.Close()

	var accountname string
	var accountname_at_cra string
	for rows.Next() {
		if err = rows.Scan(&accountname, &accountname_at_cra); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if accountname == accountname_ && accountname_at_cra == accountname_at_cra_ {
		return nil
	}

	fmt.Printf("Update Accounts set AccountName=%s, AccountNameAtCRA=%s where id=%d;\n", accountname, accountname_at_cra, accountid_)
	_, err = db.Exec("Update Accounts set AccountName=?, AccountNameAtCRA=? where id=?;", accountname, accountname_at_cra, accountid_)

	return err
}

// Delete an account and all related transactions
func DeleteAnAccount(accountname string, accountname_at_cra string) error {
	db := app.Db

	rows, err := db.Query("Select id from Accounts order by id Desc Limit 1;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var next_account_id = 0
	for rows.Next() {
		var account_id string
		if err = rows.Scan(&account_id); err != nil {
			log.Fatal(err)
		}

		max_account_id, _ := strconv.Atoi(account_id)
		next_account_id = 1 + max_account_id
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if next_account_id > 0 {
		fmt.Printf("Insert Into Accounts (id, OwnerId, AccountName, AccountNameAtCRA) Values(%d,%d,'%s','%s');\n", next_account_id, 1, accountname, accountname_at_cra)
		_, err = db.Exec("Insert Into Accounts (id, OwnerId, AccountName, AccountNameAtCRA) Values(?,?,?,?);", next_account_id, 1, accountname, accountname_at_cra)
	} else {
		err = sql.ErrNoRows
	}

	return err
}
