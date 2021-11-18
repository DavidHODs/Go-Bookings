package main

import (
	"net/http"

	"github.com/DavidHODs/bookings/config"
	"github.com/DavidHODs/bookings/handlers"
	"github.com/DavidHODs/bookings/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	chi := chi.NewRouter()

	// Middlewares
	chi.Use(middleware.Recoverer)
	chi.Use(middlewares.WriteToConsole)
	chi.Use(middlewares.NoSurf)
	chi.Use(app.Session.LoadAndSave)

	// Routes
	chi.Get("/", handlers.Repo.Home)
	chi.Get("/about", handlers.Repo.About)

	chi.Get("/generals-quarters", handlers.Repo.Generals)
	chi.Get("/majors-suites", handlers.Repo.Majors)

	chi.Get("/make-reservation", handlers.Repo.Reservation)
	chi.Post("/make-reservation", handlers.Repo.PostReservation)
	chi.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	chi.Get("/search-availability", handlers.Repo.Availability)
	chi.Post("/search-availability", handlers.Repo.PostAvailability)
	chi.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	chi.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	chi.Get("/contact", handlers.Repo.Contact)

	chi.Get("/user/login", handlers.Repo.ShowLogin)
	
	fileServer := http.FileServer(http.Dir("./static/"))
	chi.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return chi
}