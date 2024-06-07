package main

import (
	"github.com/Slava02/practiceS24/cmd/web/handlers"
	"github.com/Slava02/practiceS24/config"
	"net/http"
)

func Routes(app *config.Application) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home(app))
	mux.HandleFunc("/universe", handlers.ShowUniverse(app))
	mux.HandleFunc("/universe/create", handlers.CreateUniverse(app))

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
