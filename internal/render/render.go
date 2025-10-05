package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Chocolate529/nevarol/internal/config"
	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/justinas/nosurf"
)

var appConfig *config.AppConfig
var templatePath = "./templates/"
var functions = template.FuncMap{}

// / NewTemplates sets the appConfig.
func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using the template cache.
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template
	var err error
	if appConfig.UseChache {
		//create template chache
		templateCache = appConfig.TemplateCache
	} else {
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Println("Error creating template cache 44:", err)
			return err
		}
	}

	if len(templateCache) == 0 {
		log.Println("Error creating template cache 50:", err)
		return errors.New("no templates in cache")
	}
	//get template from cache
	curentTemplate, ok := templateCache[tmpl]
	if !ok {
		log.Println("Error creating template cache 56:", err)
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err = curentTemplate.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	//render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to response:", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	//get all files named *.page.tmpl from templates
	pages, err := filepath.Glob(fmt.Sprintf("%s*.page.tmpl", templatePath))
	if err != nil {
		log.Println("Error creating template cache 85:", err)
		return templateCache, err
	}
	if len(pages) == 0 {
		log.Println("Error creating template cache 89:", err)
		return templateCache, nil // no templates found
	}
	//range through pages
	for _, page := range pages {
		name := filepath.Base(page)

		curentTemplateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s*.layout.tmpl", templatePath))
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			curentTemplateSet, err = curentTemplateSet.ParseGlob(fmt.Sprintf("%s*.layout.tmpl", templatePath))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = curentTemplateSet
	}
	return templateCache, nil
}
