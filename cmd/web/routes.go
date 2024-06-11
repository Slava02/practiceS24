package main

import (
	"github.com/Slava02/practiceS24/cmd/web/handlers"
	"github.com/Slava02/practiceS24/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(app *config.Application) http.Handler {
	r := chi.NewRouter()

	r.Use(app.LogRequest, config.SecureHeaders)

	r.Get("/", handlers.Home(app))
	r.Get("/universe/view/{id:^[0-9]+}", handlers.ShowUniverse(app))
	r.Post("/universe/create", handlers.CreateUniverse(app))
	r.Post("/universe/create", handlers.CreateUniversePost(app))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//  TODO figure out how to manage routes in file server if we have several routes for it
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/universe/static/*", http.StripPrefix("/universe/static/", fileServer))

	return r
}
