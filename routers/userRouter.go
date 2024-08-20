package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func UserRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Get("/", controllers.GetUsers)
	r.Put("/{userId}", controllers.ChangeUserPassword)
	r.Delete("/{userId}", controllers.DeleteUser)
	r.Put("/{userId}/reset", controllers.ResetUserPassword)
}
