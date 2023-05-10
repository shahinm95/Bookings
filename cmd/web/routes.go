package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shahinm95/bookings/pkg/config"
	"github.com/shahinm95/bookings/pkg/handlers"
)

// to get rid off unused packages in mod file => go mod tidy
// to add packages => go get githubadress

func Routes(app *config.AppCongif) http.Handler {
	//handleing routes using pat package
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// the reason for using chi library is for using middlewares
	//middlewares allows us to proccess a request and perform some action on it
	mux := chi.NewRouter()

	//usnig a middleware
	//when program panics and shuts down the application , it will tell the problem
	mux.Use(middleware.Recoverer)

	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)



	return mux
}

//writing our own middleware => 
