package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
// including the template cache.
type AppConfig struct {
	UseChache bool
	TemplateCache map[string]*template.Template
	InProduction bool
	Session *scs.SessionManager
	InfoLog *log.Logger
	ErrorLog *log.Logger
}
