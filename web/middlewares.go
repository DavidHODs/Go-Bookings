package main

import (
	"net/http"

	_ "github.com/DavidHODs/bookings/config"
	"github.com/DavidHODs/bookings/helpers"

	"github.com/justinas/nosurf"
)

// WriteToConsole prints a text on the command line whenever a page is visited
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

// var app config.AppConfig

// NoSurf is a middleware that handles CSRF
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	
	return csrfHandler
}

// Auth is a middleware that handles authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

