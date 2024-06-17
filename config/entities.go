package config

import (
	"fmt"
	"github.com/Slava02/practiceS24/cmd/web/templates"
	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

const (
	ShowOnMain = 10
	ExpiesIn   = 7
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Universe       models.Storage
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
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

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, td *templates.TemplateData) {
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
