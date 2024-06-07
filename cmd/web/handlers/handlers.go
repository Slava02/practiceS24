package handlers

import (
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/config"
	"github.com/Slava02/practiceS24/pkg/models"
	"html/template"
	"net/http"
	"strconv"
)

func ShowUniverse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}
		s, err := app.Objects.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		fmt.Fprintf(w, "%+v", s)
	}
}

func CreateUniverse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
		}

		title := `Вселенная теста`
		X := 1.
		Y := 2.
		Z := 3.
		Mass := 4.
		expires := `7`

		id, err := app.Objects.Insert(title, X, Y, Z, Mass, expires)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/universe?id=%d", id), http.StatusSeeOther)
	}
}

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		files := []string{
			"./ui/html/Home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			app.ServerError(w, err)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			app.ServerError(w, err)
		}
	}
}
