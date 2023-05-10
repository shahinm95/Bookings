package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(w http.ResponseWriter , r *http.Request) {
// 		fmt.Println("you hit the page")
// 		next.ServerHTTP(w, r)
// 		http
// 	}
// }
func WriteToConsole (next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit the page")
		next.ServeHTTP(w, r)
	} )
}


// using nosurf package via middleware
// nosurf solves this problem by providing a CSRFHandler that wraps your http.Handler and checks for CSRF 
//attacks on every non-safe (non-GET/HEAD/OPTIONS/TRACE) method.
//when you build a page with form on it, you have hidden field in that form , which is long string of random 
//numbers , and they change everytime when user ges to page , so thats why we use them in middleware
func NoSurf (next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}


//Loading a session as middlerware to save and load session on every request 
func SessionLoad (next http.Handler) http.Handler {
	return sessionManagment.LoadAndSave(next)
}