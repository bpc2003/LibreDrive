package main

import (
  "log"
  "net/http"

  "libredrive/methods"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := chi.NewRouter()
  
  r.Use(middleware.Logger)

  r.Get("/api/users", methods.GetUsers)
  r.Get("/api/users/{userId}", methods.GetUserById)
  r.Put("/api/users/{userId}", methods.ChangeUserPassword)
  r.Post("/api/register", methods.CreateUser)

  log.Println("Server running on Port 8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}
