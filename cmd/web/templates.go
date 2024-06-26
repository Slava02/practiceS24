package main

import (
	"github.com/Slava02/practiceS24/pkg/models"
	"github.com/Slava02/practiceS24/pkg/validator"
	"html/template"
	"log"
	"path/filepath"
)

type TemplateData struct {
	Universe        *models.Universe
	Universes       []*models.Universe
	Form            *validator.Form
	Flash           string
	IsAuthenticated bool
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	log.Printf("Cache: %+v\n", cache)

	return cache, nil
}
