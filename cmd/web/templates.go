package main

import (
	"html/template"
	"mohamidsaiid.com/snippetbox/pkg/forms"
	"mohamidsaiid.com/snippetbox/pkg/models"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear       int
	CSRFToken         string
	Form              *forms.Form
	Flash             string
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	AuthenticatedUser *models.User
}

func humaData(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humaData,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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

	return cache, nil
}
