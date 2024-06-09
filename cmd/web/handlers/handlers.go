package handlers

import (
	"errors"
	"fmt"
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
		o, err := app.Objects.Get(id)
		log.Printf("INFO: Got Objects: %+v\n", o)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				log.Printf("ShowUniverse Server Error:\n")
				app.ServerError(w, err)
			}
			return
		}

		fmt.Fprintf(w, "%+v", o)

		//data := config.TemplateData{Object: o}
		//
		//files := []string{
		//	"./ui/html/show.page.tmpl",
		//	"./ui/html/base.layout.tmpl",
		//	"./ui/html/footer.partial.tmpl",
		//}
		//
		//ts, err := template.ParseFiles(files...)
		//if err != nil {
		//	app.ServerError(w, err)
		//	return
		//}
		//
		//err = ts.Execute(w, data)
		//if err != nil {
		//	app.ServerError(w, err)
		//}
	}
}

func CreateUniverse(app *config.Application) http.HandlerFunc {
	log.Printf("INFO: CreateUniverse API call\n")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
		}

		obj := &models.Object{
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
			},
		}

		id, err := app.Objects.Insert(obj)
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

		o, err := app.Objects.Latest(config.ShowOnMain)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		for _, obj := range o {
			fmt.Fprintf(w, "%v\n", *obj)
		}
		//files := []string{
		//	"./ui/html/Home.page.tmpl",
		//	"./ui/html/base.layout.tmpl",
		//	"./ui/html/footer.partial.tmpl",
		//}
		//
		//ts, err := template.ParseFiles(files...)
		//if err != nil {
		//	app.ErrorLog.Println(err.Error())
		//	app.ServerError(w, err)
		//	return
		//}
		//
		//err = ts.Execute(w, nil)
		//if err != nil {
		//	app.ErrorLog.Println(err.Error())
		//	app.ServerError(w, err)
		//}
	}
}
