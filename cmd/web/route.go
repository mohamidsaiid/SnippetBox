package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)
func (app *Application) route() http.Handler{
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	return mux
}