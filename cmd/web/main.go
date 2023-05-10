package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shahinm95/bookings/pkg/config"
	"github.com/shahinm95/bookings/pkg/handlers"
	"github.com/shahinm95/bookings/pkg/render"
)

// for running all go files in directory => go run .


// making template cache here , defining outside of main func to give access other files in main package like middlerware
var app config.AppCongif

// defining session outside of main func to give access other files in main package, like handlers with help of config
var sessionManagment *scs.SessionManager


func main () {


	// change this in InProduction mode to true
	app.InProduction = false

	templateCache , err :=render.CreateTemplateCacheThird()
	if err != nil {log.Fatal(err)}
	app.TemplateCache = templateCache
	app.UseCache = false
	//calling a function in render to send templates data to render file
	render.NewTemplates(&app) // we refrence to app because in function we used pointer as type

	//calling functions in handlers to give access to config in handlers file
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// handling request to specific address and responding to it
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	// adding session managment package to main 
	// by default session is using cookies to store
	sessionManagment = scs.New()
	app.SessionManager = sessionManagment
	sessionManagment.Lifetime = 30 * time.Minute // 30 minutes session lifetime
	// since we use cookies we need to set some aparameters for cookies
	sessionManagment.Cookie.Persist = true // should session presisit after closing the window
	sessionManagment.Cookie.SameSite = http.SameSiteLaxMode
	sessionManagment.Cookie.Secure = app.InProduction // this will insist cookies will be incrypted and uses https protocol



	//satrting a server and listening for requests
	// err = http.ListenAndServe(":3000", nil)
	// if err != nil {fmt.Println("error listing server", err)}
	fmt.Println("listening to requests on Port 3000")

	//handing routes using render.go file using pat package
	//handling request to specific address with another way via contruction variable via http.Server struct
	serv := http.Server{
		Addr: ":3000",
		Handler: Routes(&app),
	}
	err = serv.ListenAndServe()
	if err != nil {log.Fatal(err)}

}



// for installing session managment package => go get github.com/alexedwards/scs/v2
