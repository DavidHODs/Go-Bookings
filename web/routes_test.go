package main

import (
	"fmt"
	"testing"

	"github.com/DavidHODs/bookings/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type){
	case *chi.Mux:

	default:
		t.Error(fmt.Sprintf("type is not http.Handler but %T",v))
	}
}