package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	"github.com/rainclear/accroo/pkg/config"
	"github.com/rainclear/accroo/pkg/handlers"
	"github.com/rainclear/accroo/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig

func main() {
	// change this to true when in production
	app.InProduction = false
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println("Starting application on port: ", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
