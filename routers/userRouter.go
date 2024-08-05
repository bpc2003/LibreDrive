package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func GroupRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Use(customMiddleware.IsAuth)
	r.Get("/", controllers.GetUsers)
}

func IndividualRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Use(customMiddleware.IsAuth)
	r.Get("/", controllers.GetUserById)
	r.Put("/", controllers.ChangeUserPassword)
	r.Delete("/", controllers.DeleteUser)
}
