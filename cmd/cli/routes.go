package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	auth := alice.New(app.LoginMiddleware)

	//home related routes
	router.HandlerFunc(http.MethodGet, "/", app.Home)
	//books related routes
	router.HandlerFunc(http.MethodGet, "/book/listing", app.BookListing)        // all the book listing
	router.HandlerFunc(http.MethodGet, "/book/search/:isbn/", app.BookInfo)     // review of given isbn
	router.Handler(http.MethodPost, "/book/create", auth.ThenFunc(app.AddBook)) // create a new book info different isbn
	//review related routes
	router.HandlerFunc(http.MethodGet, "/review/listing", app.ReviewListing)              // all the reviews
	router.Handler(http.MethodGet, "/myreview/", auth.ThenFunc(app.MyReview))             // review of logged in user
	router.HandlerFunc(http.MethodGet, "/review/search/:isbn/", app.ReviewSearch)         // review of given isbn
	router.Handler(http.MethodPost, "/review/create", auth.ThenFunc(app.AddReview))       // create review
	router.Handler(http.MethodGet, "/review/delete/:id", auth.ThenFunc(app.DeleteReview)) // delete your own review
	//user related routes
	router.HandlerFunc(http.MethodPost, "/user/forget_password/", app.ForgetPasswordPost) // to create forget password request
	router.HandlerFunc(http.MethodPost, "/user/register", app.UserRegister)               // to register
	router.HandlerFunc(http.MethodPost, "/user/new_password/:uri", app.NewPasswordPost)   //after forget password req uri created to change passw
	router.HandlerFunc(http.MethodGet, "/user/activation/:uri", app.UserActivation)       // after registration uri created authentication
	router.HandlerFunc(http.MethodPost, "/user/login", app.UserLogin)                     // login
	router.Handler(http.MethodPost, "/user/logout", auth.ThenFunc(app.UserLogout))        // logout
	standard := alice.New(app.logRequest, secureHeaders)
	return standard.Then(router)
}
