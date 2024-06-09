package main

import (
	"github.com/Slava02/practiceS24/cmd/web/handlers"
	"github.com/Slava02/practiceS24/config"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func Routes(app *config.Application) http.Handler {
	r := chi.NewRouter()

	r.Use(app.LogRequest, config.SecureHeaders)

	r.Get("/", handlers.Home(app))
	r.Get("/universe/view/{id:^[0-9]+}", handlers.ShowUniverse(app))
	r.Post("/universe/create", handlers.CreateUniverse(app))
	r.Post("/universe/create", handlers.CreateUniversePost(app))

	filesDir := http.Dir("./ui/static/")
	FileServer(r, "/static", filesDir)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
