package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func NavRoute(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Get("/", controllers.Nav)
}
