package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	dynamicMiddleware := alice.New(app.sessions.Enable, noSurf, app.authenticate)
	mux := pat.New()
	// mux := http.NewServeMux()
	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.showSnippet)))

	// 	GET /user/signup signupUserForm Display the user signup form
	// POST /user/signup signupUser Create a new user
	// GET /user/login loginUserForm Display the user login form
	// POST /user/login loginUser Authenticate and login the user
	// POST /user/logout logoutUser Logout the user

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.logoutUser)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	return standardMiddleware.Then(mux)
}
