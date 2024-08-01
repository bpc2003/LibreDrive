package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"libredrive/controllers"
	"libredrive/customMiddleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Use(customMiddleware.Auth)
		r.Get("/", controllers.GetFiles)
		r.Post("/", controllers.UploadFile)
	})

	r.Route("/api/users", func(r chi.Router) {
		r.Use(customMiddleware.Auth)
		r.Use(customMiddleware.IsAuth)
		r.Get("/", controllers.GetUsers)
	})

	r.Route("/api/users/{userId}", func(r chi.Router) {
		r.Use(customMiddleware.Auth)
		r.Use(customMiddleware.IsAuth)
		r.Get("/", controllers.GetUserById)
		r.Put("/", controllers.ChangeUserPassword)
		r.Delete("/", controllers.DeleteUser)
	})

	r.Route("/api/register", func(r chi.Router) {
		r.Use(customMiddleware.Auth)
		r.Use(customMiddleware.IsAuth)
		r.Post("/", controllers.CreateUser)
	})

	r.Post("/api/login", controllers.LoginUser)

	log.Println("Server running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
