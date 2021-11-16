package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DavidHODs/bookings/config"
	"github.com/DavidHODs/bookings/models"
	"github.com/DavidHODs/bookings/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// data stored in sesion
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog  

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	chi := chi.NewRouter()

	// Middlewares
	chi.Use(middleware.Recoverer)
	chi.Use(WriteToConsole)
	// chi.Use(NoSurf)
	chi.Use(app.Session.LoadAndSave)

	// Routes
	chi.Get("/", Repo.Home)
	chi.Get("/about", Repo.About)

	chi.Get("/generals-quarters", Repo.Generals)
	chi.Get("/majors-suites", Repo.Majors)

	chi.Get("/make-reservation", Repo.Reservation)
	chi.Post("/make-reservation", Repo.PostReservation)
	chi.Get("/reservation-summary", Repo.ReservationSummary)

	chi.Get("/search-availability", Repo.Availability)
	chi.Post("/search-availability", Repo.PostAvailability)
	chi.Post("/search-availability-json", Repo.AvailabilityJSON)

	chi.Get("/contact", Repo.Contact)
	
	fileServer := http.FileServer(http.Dir("./static/"))
	chi.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return chi
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*page.gohtml", pathToTemplate))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
		if err != nil {
		return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
			if err != nil {
			return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

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