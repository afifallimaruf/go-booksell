package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	dynamicMiddleware := alice.New(app.session.Enable)

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/", dynamicMiddleware.Then(http.HandlerFunc(app.indexHandler))).Methods("GET")
	router.Handle("/books", dynamicMiddleware.Then(http.HandlerFunc(app.booksHandler))).Methods("GET")
	router.Handle("/about", dynamicMiddleware.Then(http.HandlerFunc(app.aboutHandler))).Methods("GET")

	// Add Book
	router.Handle("/add-books", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.addBooksForm))).Methods("GET")
	router.Handle("/add-books", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.addBooks))).Methods("POST")

	// Add to chart
	router.Handle("/add-chart", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.addChart))).Methods("GET")

	// Category
	// non-fiction
	router.Handle("/category/non-fiction", dynamicMiddleware.Then(http.HandlerFunc(app.nonFictionHandler))).Methods("GET")

	// Signup
	router.Handle("/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupPage))).Methods("GET")
	router.Handle("/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupUser))).Methods("POST")

	// Login
	router.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginPage))).Methods("GET")
	router.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginHandler))).Methods("POST")

	// Logout
	router.Handle("/logout", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.logoutHandler)))

	fs := http.FileServer(http.Dir("views/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return app.logRequest(router)
}
