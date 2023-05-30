package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shahinm95/bookings/internal/config"
)

func TestRoutes(t *testing.T) {

	var app config.AppConfig
	rt := routes(&app)

	switch re := rt.(type) {
	case *chi.Mux:
	// do nothing
	default:
		t.Errorf("type is not *chi.Mux, type is %T", re)
	}
}
