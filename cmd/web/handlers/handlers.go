package handlers

import (
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/cmd/web/templates"
	"github.com/Slava02/practiceS24/config"
	"github.com/Slava02/practiceS24/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func ShowUniverse(app *config.Application) http.HandlerFunc {
	log.Printf("\nShowUniverse API call\n")
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}
		universe, err := app.Universe.Get(id)
		log.Printf("INFO: Got Universe: %+v\n", universe)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				log.Printf("ShowUniverse Server Error:\n")
				app.ServerError(w, err)
			}
			return
		}

		//fmt.Fprintf(w, "%+v", universe)
		app.Render(w, r, "show.page.tmpl", &templates.TemplateData{
			Universe: universe,
		})
	}
}

func CreateUniverse(app *config.Application) http.HandlerFunc {
	log.Printf("INFO: CreateUniverse API call\n")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
		}

		obj := &models.Universe{
			Title: `Вселенная теста`,
			Params: []*models.Params{
				{
					Coord: &models.Coord{
						X: 1.,
						Y: 2.,
						Z: 3.,
					},
					Mass: 4.,
				},
				{
					Coord: &models.Coord{
						X: 10.,
						Y: 20.,
						Z: 30.,
					},
					Mass: 40.,
				},
			},
		}

		id, err := app.Universe.Insert(obj)
		if err != nil {
			log.Printf("ERROR: CreateUniverse ServerError:\n")
			app.ServerError(w, err)
			return
		}

		log.Printf("INFO: Redirect to /universe?id=%d\n", id)
		http.Redirect(w, r, fmt.Sprintf("/universe?id=%d", id), http.StatusSeeOther)
	}
}

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		universes, err := app.Universe.Latest(config.ShowOnMain)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		//for _, obj := range universe {
		//	fmt.Fprintf(w, "%v\n", *obj)
		//}

		app.Render(w, r, "home.page.tmpl", &templates.TemplateData{
			Universes: universes,
		})
	}
}
