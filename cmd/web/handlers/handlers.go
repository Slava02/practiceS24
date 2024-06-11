package handlers

import (
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
	log.Printf("INFO: I AM IN CreateUniversePost HANDLER")
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			log.Printf("Couldnt parse form: %x\n", err)
			app.ClientError(w, http.StatusBadRequest)
		}

		x_s := r.PostForm.Get("x")
		y_s := r.PostForm.Get("y")
		z_s := r.PostForm.Get("z")
		m_s := r.PostForm.Get("mass")

		log.Printf("PARSED: X-%s, Y-%s, Z-%s, M-%s\n", x_s, y_s, z_s, m_s)

		x, err := strconv.ParseFloat(x_s, 64)
		y, err := strconv.ParseFloat(y_s, 64)
		z, err := strconv.ParseFloat(z_s, 64)
		mass, err := strconv.ParseFloat(m_s, 64)

		if err != nil {
			app.ServerError(w, err)
		}

		p := models.NewParams(x, y, z, mass)

		title := r.PostForm.Get("title")
		expiers, err := strconv.Atoi(r.PostForm.Get("expires"))
		params := []*models.Params{p}

		obj := models.NewUniverse(title, params, expiers)

		id, err := app.Universe.Insert(obj)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/universe/view/%d", id), http.StatusSeeOther)
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
