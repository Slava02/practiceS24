package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(app *Application) http.Handler {
	r := chi.NewRouter()

	r.Use(app.SessionManager.LoadAndSave, app.LogRequest, SecureHeaders)

	r.Get("/", Home(app))
	r.Get("/universe/view/{id:^[0-9]+}", ShowUniverse(app))
	r.Get("/universe/create", CreateUniverse(app))
	r.Post("/universe/create", CreateUniversePost(app))

	r.Get("/user/signup", UserSignup(app))
	r.Post("/user/signup", UserSignupPost(app))
	r.Get("/user/login", UserLogin(app))
	r.Post("/user/login", UserLoginPost(app))
	r.Post("/user/logout", UserLogoutPost(app))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//  TODO figure out how to manage routes in file server if we have several routes for it
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/universe/static/*", http.StripPrefix("/universe/static/", fileServer))

	return r
}
