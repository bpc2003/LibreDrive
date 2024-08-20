package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func GroupRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Get("/", controllers.GetUsers)
}

func IndividualRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Put("/", controllers.ChangeUserPassword)
	r.Delete("/", controllers.DeleteUser)
	r.Put("/reset", controllers.ResetUserPassword)
}
