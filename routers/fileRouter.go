package routers

import (
	"github.com/go-chi/chi/v5"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func FileRoutes(r chi.Router) {
	r.Use(customMiddleware.Auth)
	r.Get("/", controllers.GetFiles)
	r.Post("/", controllers.UploadFile)
 	r.Get("/{fileName}", controllers.GetFile)
	r.Delete("/{fileName}", controllers.DeleteFile)
}
