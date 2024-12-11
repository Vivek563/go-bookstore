package main

import (
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ldata")
		if err != nil || cookie.Value == "" || len(cookie.Value) != 40 {
			app.notFound(w)
			app.infoLog.Print("Invalid Logout")
			return
		}

		userID, err := app.models.Users.ValidUser(cookie.Value)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if userID > 0 {
			app.isAuthenticated = true
		}

		if !app.isAuthenticated {
			http.Redirect(w, r, "/api/user/login/", http.StatusSeeOther)
			return
		}

		app.user_id = userID
		next.ServeHTTP(w, r)

	})
}
