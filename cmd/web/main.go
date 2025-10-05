package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chocolate529/nevarol/internal/config"
	"github.com/Chocolate529/nevarol/internal/handlers"
	"github.com/Chocolate529/nevarol/internal/helpers"
	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/Chocolate529/nevarol/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal("Failed to run setup")
	}
	
	fmt.Printf("Starting app on port %s", portNumber)
	
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func run() error {
	//set the value type that is stored in the session
	gob.Register(models.Reservation{})
	///change to true when secure connection
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction //because 8080 connection is not secure

	appConfig.Session = session

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateChache, err := render.CreateTemplateCache()
	if err != nil {
		return err
	}
	appConfig.TemplateCache = templateChache
	appConfig.UseChache = false

	repo := handlers.NewRepo(&appConfig)

	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)

	return nil
}
