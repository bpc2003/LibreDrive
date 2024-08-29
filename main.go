// LibreDrive - A privacy oriented version of Google Drive
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libredrive/controllers"
	"libredrive/global"
	"libredrive/routers"
)

func main() {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.Use(middleware.Logger)

	r.Handle("/*", fs)
	r.Route("/nav", routers.NavRoute)
	r.Route("/api/files", routers.FileRoutes)
	r.Route("/api/users", routers.UserRoutes)
	r.Route("/api/register", routers.RegisterRoute)
	r.Post("/api/login", controllers.LoginUser)

	log.Println("Server running on Port " + global.PORT)
	log.Fatal(http.ListenAndServe(":" + global.PORT, r))
}
