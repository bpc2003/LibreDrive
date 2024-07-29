package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libredrive/methods"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/api/users", methods.GetUsers)
	r.Get("/api/users/{userId}", methods.GetUserById)
	r.Put("/api/users/{userId}", methods.ChangeUserPassword)
	r.Delete("/api/users/{userId}", methods.DeleteUser)

	r.Post("/api/register", methods.CreateUser)

	r.Post("/api/login", methods.LoginUser)

	log.Println("Server running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
