package main

import (
	"fmt"
	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

const (
	ShowOnMain = 10
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Universe       models.UniverseModel
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
	Users          models.UserModel
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.ServerError(w, err)
	}
}

func (app *Application) IsAuthenticated(r *http.Request) bool {
	return app.SessionManager.Exists(r.Context(), "authenticatedUserID")
}
