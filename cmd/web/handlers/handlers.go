package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/cmd/web"
	"github.com/Slava02/practiceS24/cmd/web/templates"
	"github.com/Slava02/practiceS24/pkg/forms"
	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/Slava02/practiceS24/pkg/validator"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func ShowUniverse(app *main.Application) http.HandlerFunc {
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
			Flash:    app.SessionManager.PopString(r.Context(), "flash"),
		})
	}
}

func CreateUniverse(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "create.page.tmpl", new(templates.TemplateData))
	}
}

func CreateUniversePost(app *main.Application) http.HandlerFunc {
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

		obj.CheckField(validator.NotBlank(obj.Title), "title", "Это поле не может быть пустым")
		obj.CheckField(validator.MaxChars(obj.Title, 100), "title", "Это поле не может быть больше 100 символов")
		obj.CheckField(validator.NotNill(obj.Params), "params", "Параметры не могут быть пустыми")
		obj.CheckField(validator.PermittedInt(obj.ExpiresIn, 1, 7, 365), "expires", "Это поле должно равняться 1, 7 или 365")

		log.Printf("DECODED: %+v PARAMS: %+v", obj, obj.Params)

		if !obj.Valid() {
			app.ClientError(w, http.StatusUnprocessableEntity)
		} else {
			id, err := app.Universe.Insert(&obj)
			if err != nil {
				app.ServerError(w, err)
				return
			}
			app.SessionManager.Put(r.Context(), "flash", "Вселенная успешно создана!")
			fmt.Fprintf(w, fmt.Sprintf("%d", id))
		}
	}
}

func Home(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		universes, err := app.Universe.Latest(main.ShowOnMain)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.Render(w, r, "home.page.tmpl", &templates.TemplateData{
			Universes: universes,
		})
	}
}

func UserSignup(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "home.page.tmpl", &templates.TemplateData{
			Forms: forms.NewForm(nil),
		})
		//fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
	}
}

func UserSignupPost(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create a new user...")
	}
}

func UserLogin(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Display a HTML form for logging in a user...")
	}
}

func UserLoginPost(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Display a HTML form for logging in a user...")
	}
}

func UserLogoutPost(app *main.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Display a HTML form for logging in a user...")
	}
}
