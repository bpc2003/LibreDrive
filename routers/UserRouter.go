package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/middleware"
)

// UserRoutes - various user routes
func UserRoutes(r chi.Router) {
	r.Use(middleware.Auth)
	r.Get("/", controllers.GetUsers)
	r.Put("/{userId}", controllers.ChangeUserPassword)
	r.Delete("/{userId}", controllers.DeleteUser)
	r.Put("/{userId}/reset", controllers.ResetUserPassword)
}
