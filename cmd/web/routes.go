// this file contains all the serveMux routes and its own handlers
package main


import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/bmizerany/pat"
)


func (app *Application) routes() http.Handler{
	

	standardMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleWare := alice.New(app.session.Enable, app.noSurf, app.authenticate)

	/*	// initalizing my own serveMux
	mux := http.NewServeMux()
	// some routes to specific handelres from the type fixed root
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	
	// initsating fileSever with the http package this returns and http.Handler to be used in the routing
	fileServer := http.FileServer(http.Dir("./ui/static"))
	
	// this takes treeroot route and its own http.Handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	// return your own serveMux*/


	mux := pat.New()
	mux.Get("/", dynamicMiddleWare.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleWare.Append(app.requiredAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleWare.Append(app.requiredAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleWare.ThenFunc(app.showSnippet))
	

	mux.Get("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUserForm))	
	mux.Post("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUser))	
	mux.Get("/user/login", dynamicMiddleWare.Append(app.alreadyAuthenticated).ThenFunc(app.loginUserForm))	
	mux.Post("/user/login", dynamicMiddleWare.Append(app.alreadyAuthenticated).ThenFunc(app.loginUser))	
	mux.Post("/user/logout", dynamicMiddleWare.Append(app.requiredAuthenticatedUser).ThenFunc(app.logoutUser))	

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))



	return standardMiddleWare.Then(mux)
}
