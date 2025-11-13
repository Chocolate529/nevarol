package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chocolate529/nevarol/internal/config"
	"github.com/Chocolate529/nevarol/internal/driver"
	"github.com/Chocolate529/nevarol/internal/handlers"
	"github.com/Chocolate529/nevarol/internal/helpers"
	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/Chocolate529/nevarol/internal/render"
	"github.com/Chocolate529/nevarol/internal/repository"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal("Failed to run setup:", err)
	}
	defer db.Pool.Close()
	
	fmt.Printf("Starting app on port %s\n", portNumber)
	
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func run() (*driver.DB, error) {
	//set the value type that is stored in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	
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

	// Database connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "nevarol"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	appConfig.InfoLog.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		appConfig.ErrorLog.Println("Cannot connect to database! Dying...")
		return nil, err
	}

	// Run migrations
	err = db.RunMigrations()
	if err != nil {
		appConfig.ErrorLog.Println("Cannot run migrations:", err)
		return nil, err
	}

	// Setup database repository
	appConfig.DB = repository.NewDatabaseRepo(db.Pool)

	templateChache, err := render.CreateTemplateCache()
	if err != nil {
		return nil, err
	}
	appConfig.TemplateCache = templateChache
	appConfig.UseChache = false

	repo := handlers.NewRepo(&appConfig)

	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)

	return db, nil
}
