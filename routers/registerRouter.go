package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func RegisterRoute(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Use(customMiddleware.IsAuth)
	r.Post("/", controllers.CreateUser)
}
