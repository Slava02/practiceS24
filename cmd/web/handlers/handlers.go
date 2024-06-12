package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/cmd/web/templates"
	"github.com/Slava02/practiceS24/config"
	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func ShowUniverse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParamFromCtx(r.Context(), "id"))
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}

		universe, err := app.Universe.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		app.Render(w, r, "show.page.tmpl", &templates.TemplateData{
			Universe: universe,
		})
	}
}

func CreateUniverse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "create.page.tmpl", new(templates.TemplateData))
	}
}

func CreateUniversePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Couldnt parse form: %x\n", err)
			app.ClientError(w, http.StatusBadRequest)
		}

		var obj models.Universe

		err = json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := app.Universe.Insert(&obj)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		//http.Redirect(w, r, fmt.Sprintf("/universe/view/%d", id), http.StatusSeeOther)
		fmt.Fprintf(w, fmt.Sprintf("%d", id))
	}
}

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		universes, err := app.Universe.Latest(config.ShowOnMain)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.Render(w, r, "home.page.tmpl", &templates.TemplateData{
			Universes: universes,
		})
	}
}
