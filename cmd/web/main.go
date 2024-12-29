package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	"github.com/rainclear/troomate/pkg/config"
	"github.com/rainclear/troomate/pkg/dbm"
	"github.com/rainclear/troomate/pkg/handlers"
	"github.com/rainclear/troomate/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig

func main() {
	// change this to true when in production
	app.InProduction = false
	app.DBPath = "profiles/sandprofile.db"
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	dbm.NewDbm(&app)

	err := dbm.OpenDb()
	defer dbm.CloseDb()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.Owner)
	fmt.Println("Number of Account: ", len(app.Accounts))
	for _, account := range app.Accounts {
        fmt.Println(account)
    }

	log.Println("Starting application on port: ", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
