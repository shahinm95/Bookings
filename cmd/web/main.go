package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shahinm95/bookings/internal/config"
	"github.com/shahinm95/bookings/internal/driver"
	"github.com/shahinm95/bookings/internal/handlers"
	"github.com/shahinm95/bookings/internal/helpers"
	"github.com/shahinm95/bookings/internal/models"
	"github.com/shahinm95/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	go listenForMail()
	fmt.Println("starting mail listener...")

	fmt.Printf("Staring application on port %s", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	// read flags
	inProduction := flag.Bool("production", true, "application in production mode")
	useCache := flag.Bool("cache", true , "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "","Database User name")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable , prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser =="" {
		fmt.Println("Please specify required flags" )
		os.Exit(1)
	}
	if *dbPass == "" {
		fmt.Println("Please specify required flags" )
		os.Exit(1)
	}
	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache


	// error handling
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	//what I'm going to put in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	//connect to my database
	log.Println("connecting to database...")
	conectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(conectionString)
	if err != nil {
		log.Fatal("error connecting to database", err)
	}
	log.Println("connected to database successfully")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}
