package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"libredrive/methods"
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

	r.Route("/api/users", func(r chi.Router) {
		r.Use(customMiddleware.Auth)
		r.Get("/", methods.GetUsers)
		r.Get("/{userId}", methods.GetUserById)
		r.Put("/{userId}", methods.ChangeUserPassword)
		r.Delete("/{userId}", methods.DeleteUser)
	})

	r.Post("/api/register", methods.CreateUser)

	r.Post("/api/login", methods.LoginUser)

	log.Println("Server running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
