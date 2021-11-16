package config

import (
	"html/template"
	"log"

	"github.com/DavidHODs/bookings/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
type AppConfig struct {
	UseCache 		bool
	TemplateCache 	map[string]*template.Template
	ErrorLog		*log.Logger
	InfoLog 		*log.Logger
	InProduction 	bool
	Session 		*scs.SessionManager
	MailChan		chan models.MailData
}