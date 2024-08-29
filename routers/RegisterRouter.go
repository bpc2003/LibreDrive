package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/middleware"
)

// RegisterRoute - route for registering users
func RegisterRoute(r chi.Router) {
	r.Use(middleware.Auth)
	r.Post("/", controllers.CreateUser)
}
