package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/middleware"
)

func NavRoute(r chi.Router) {
	r.Use(middleware.Auth)
	r.Get("/", controllers.Nav)
}
