package config

import (
	"html/template"
	"log"

	"github.com/Chocolate529/nevarol/internal/email"
	"github.com/Chocolate529/nevarol/internal/repository"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
// including the template cache.
type AppConfig struct {
	UseChache     bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DB            *repository.DatabaseRepo
	EmailConfig   *email.Config
}
