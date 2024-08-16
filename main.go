package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"libredrive/controllers"
	"libredrive/routers"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.Use(middleware.Logger)

	r.Handle("/", fs)
	r.Handle("/login/", fs)
	r.Route("/api/files", routers.FileRoutes)
	r.Route("/api/users", routers.GroupRoutes)
	r.Route("/api/users/{userId}", routers.IndividualRoutes)
	r.Route("/api/register", routers.RegisterRoute)
	r.Route("/nav", routers.NavRoute)
	r.Post("/api/login", controllers.LoginUser)

	log.Println("Server running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
