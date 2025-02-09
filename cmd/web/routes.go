package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rainclear/troomate/pkg/config"
	"github.com/rainclear/troomate/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/accounts", handlers.Repo.Accounts)
	mux.Get("/new_account", handlers.Repo.NewAccount)
	mux.Get("/modify_account", handlers.Repo.ModifyAccount)

	mux.Post("/new_account", handlers.Repo.PostNewAccount)
	mux.Post("/delete_an_account", handlers.Repo.DeleteAnAccount)
	mux.Post("/edit_an_account", handlers.Repo.EditAnAccount)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
