package handlers

import (
	
	"net/http"

	"github.com/Chocolate529/nevarol/internal/config"
	
	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/Chocolate529/nevarol/internal/render"
)

// Repository holds the application configuration.
// It is used to pass the app configuration to the handlers.
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository with the given appConfig.
// It initializes the repository with the app configuration.
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers.
func NewHandlers(r *Repository) {
	Repo = r
}

var Repo *Repository

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Store(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "store.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Shipping(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "shipping.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Checkout(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "checkout.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Account(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "account.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{})
}
