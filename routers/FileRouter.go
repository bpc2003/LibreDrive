// routers - custom routers
package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/middleware"
)

// FileRouter - routes for user files
func FileRoutes(r chi.Router) {
	r.Use(middleware.Auth)
	r.Get("/", controllers.GetFiles)
	r.Post("/", controllers.UploadFile)
	r.Get("/{fileName}", controllers.GetFile)
	r.Delete("/{fileName}", controllers.DeleteFile)
}
