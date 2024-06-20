package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/Slava02/practiceS24/pkg/validator"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func ShowUniverse(app *Application) http.HandlerFunc {
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

		app.Render(w, r, "show.page.tmpl", &TemplateData{
			Universe:        universe,
			Flash:           app.SessionManager.PopString(r.Context(), "flash"),
			IsAuthenticated: app.IsAuthenticated(r),
		})
	}
}

func CreateUniverse(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "create.page.tmpl", &TemplateData{
			Flash:           app.SessionManager.PopString(r.Context(), "flash"),
			IsAuthenticated: app.IsAuthenticated(r),
		})
	}
}

func CreateUniversePost(app *Application) http.HandlerFunc {
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

func Home(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		universes, err := app.Universe.Latest(ShowOnMain)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.Render(w, r, "home.page.tmpl", &TemplateData{
			Universes:       universes,
			Flash:           app.SessionManager.PopString(r.Context(), "flash"),
			IsAuthenticated: app.IsAuthenticated(r),
		})
	}
}

func UserSignupPost(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form := validator.NewForm(r.PostForm)
		form.Required("name", "email", "password")
		form.MatchesPattern("email", validator.EmailRX)
		form.MinLength("password", 5)

		err = app.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err == models.ErrDuplicateEmail {
			form.Errors.Add("email", "Почта уже используется")
		} else if err != nil {
			app.ServerError(w, err)
			return
		}

		if !form.Valid() {
			app.Render(w, r, "signup.page.tmpl", &TemplateData{
				Form:            form,
				IsAuthenticated: app.IsAuthenticated(r),
			})
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Вы успешно зарегестрировались. Войдите")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func UserSignup(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "signup.page.tmpl", &TemplateData{
			Form:            validator.NewForm(nil),
			IsAuthenticated: app.IsAuthenticated(r),
		})
	}
}

func UserLoginPost(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form := validator.NewForm(r.PostForm)
		form.Required("email", "password")

		if !form.Valid() {
			app.Render(w, r, "login.page.tmpl", &TemplateData{
				Form:            form,
				Flash:           app.SessionManager.PopString(r.Context(), "flash"),
				IsAuthenticated: app.IsAuthenticated(r),
			})
			return
		}
		id, err := app.Users.Authenticate(form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.Add("nonFieldError", "Неккоректные email или пароль")
				app.Render(w, r, "login.page.tmpl", &TemplateData{
					Form:            form,
					Flash:           app.SessionManager.PopString(r.Context(), "flash"),
					IsAuthenticated: app.IsAuthenticated(r),
				})
			} else {
				app.ServerError(w, err)
			}
			return
		}
		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, err)
			return
		}
		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
		http.Redirect(w, r, "/universe/create", http.StatusSeeOther)
	}
}

func UserLogin(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, "login.page.tmpl", &TemplateData{
			Form:            validator.NewForm(nil),
			Flash:           app.SessionManager.PopString(r.Context(), "flash"),
			IsAuthenticated: app.IsAuthenticated(r),
		})
	}
}

func UserLogoutPost(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			return
		}

		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		app.SessionManager.Put(r.Context(), "flash", "Вы успешно вышли!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
