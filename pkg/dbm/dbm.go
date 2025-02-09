package dbm

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/rainclear/troomate/pkg/config"
	"github.com/rainclear/troomate/pkg/models"
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

	fi, err := os.Stat(app.DBPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s size: %v\n", app.DBPath, fi.Size())
	return err
}

func GetAccountInfo(accountId int64) (models.Account, error) {
	var account models.Account
	db := app.Db

	// Fetch account details for the given ID
	//	row := db.QueryRow("SELECT id, AccountName, Institution, AccountNumber, AccountNameAtCRA, AccountType, AccountPurpose FROM accounts WHERE id = ?", accountId)
	row := db.QueryRow("SELECT id, AccountName, AccountNameAtCRA FROM accounts WHERE id = ?", accountId)
	err := row.Scan(&account.ID, &account.AccountName, &account.AccountNameAtCRA)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No record found")
		} else {
			fmt.Println("Query error:", err)
		}
	}

	return account, err
}

func ListAccounts() ([]models.Account, error) {
	var accounts []models.Account

	db := app.Db

	rows, err := db.Query("Select id, AccountName, AccountNameAtCRA from Accounts order by AccountName;")
	if err != nil {
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		if err = rows.Scan(&account.ID, &account.AccountName, &account.AccountNameAtCRA); err != nil {
			return accounts, err
		}

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

	var max_account_id int64
	var next_account_id int64
	max_account_id = 0
	next_account_id = 0
	for rows.Next() {
		if err = rows.Scan(&max_account_id); err != nil {
			return err
		}
	}
	next_account_id = 1 + max_account_id

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

// Delete an account and all related transactions
func DeleteAnAccount(accountId int64) error {
	db := app.Db

	stmt, err := db.Prepare("Delete from Accounts Where id = ?;")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(accountId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Number of rows deleted: %d\n", rowsAffected)
	return err
}

func UpdateAnAccount(accountid_ int64, accountname_ string, accountname_at_cra_ string) error {
	db := app.Db

	fmt.Printf("Update Accounts set AccountName=%s, AccountNameAtCRA=%s where id=%d;\n", accountname_, accountname_at_cra_, accountid_)

	stmt, err := db.Prepare("Update Accounts set AccountName=?, AccountNameAtCRA=? where id=?;")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	rowsAffected, err := stmt.Exec(accountname_, accountname_at_cra_, accountid_)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Number of rows updated: %d\n", rowsAffected)
	return err
}
