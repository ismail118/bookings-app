package main

import (
	"github.com/go-chi/chi"
	"github.com/ismail118/bookings-app/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("%T not type *chi.Mux", v)
	}
}
