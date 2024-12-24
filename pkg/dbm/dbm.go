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

	rows, err := db.Query("Select * from Accounts Order by id;")
	if err != nil {
		return err
	}

	for rows.Next() {
		var account string
		if err = rows.Scan(&account); err != nil {
			return err
		}
		app.Accounts = append(app.Accounts, account)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	fi, err := os.Stat(app.DBPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s size: %v\n", app.DBPath, fi.Size())
	return nil
}

// func ListAccountCategories() ([]string, error) {
// 	var account_categories []string
// 	db := app.Db

// 	rows, err := db.Query("Select AccountCategory from AccountCategories Order by id;")
// 	if err != nil {
// 		return account_categories, err
// 	}

// 	for rows.Next() {
// 		var account_category string
// 		if err = rows.Scan(&account_category); err != nil {
// 			return account_categories, err
// 		}
// 		account_categories = append(account_categories, account_category)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return account_categories, err
// 	}

// 	return account_categories, nil
// }
