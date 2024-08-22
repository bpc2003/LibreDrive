package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func RegisterRoute(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Post("/", controllers.CreateUser)
}
