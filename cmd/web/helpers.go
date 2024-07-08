// this file for some helper methods
package main

import (
	"time"
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"github.com/justinas/nosurf"
	"mohamidsaiid.com/snippetbox/pkg/models"
)


// serverError that prints the error and its own stack
// the stack that returns contain each error and what it refers to
// make an error and you would see what I mean
func (app *Application) severError(w http.ResponseWriter, err error) {
	// trance is string type which contians the err and its own stack
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


// clientError appears when some request from the client side fail to get what it need as response
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}


// not found error occurs on the client side as it sends bad requests
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}


func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.severError(w,  fmt.Errorf("the template %s does not exist", name))
		return
	}
	
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil{
		app.severError(w, err)
		return
	}
	buf.WriteTo(w)
}


func (app *Application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	return td
}


func (app *Application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	} 

	return user
}