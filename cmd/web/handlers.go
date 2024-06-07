package web

import (
	"fmt"
	"net/http"
	"strconv"
)

func ShowUniverse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}
	fmt.Fprintf(w, "Отображение вселенной c id = %d", id)
}

func CreateUniverse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method is not allowed", 405)
	}
	fmt.Fprintf(w, "Форма для создания новой вселенной")
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Привет из мультивселеной!")
}
