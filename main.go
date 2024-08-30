// LibreDrive - A privacy oriented version of Google Drive
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libredrive/controllers"
	"libredrive/global"
	"libredrive/routers"
)

//go:embed public/*
var content embed.FS

func main() {
	r := chi.NewRouter()
	public, _ := fs.Sub(content, "public")
	fs := http.FileServer(http.FS(public))

	r.Use(middleware.Logger)

	r.Handle("/*", fs)
	r.Route("/nav", routers.NavRoute)
	r.Route("/api/files", routers.FileRoutes)
	r.Route("/api/users", routers.UserRoutes)
	r.Post("/api/login", controllers.LoginUser)
	r.Post("/api/register", controllers.CreateUser)

	log.Println("Server running on Port " + global.PORT)
	log.Fatal(http.ListenAndServe(":" + global.PORT, r))
}
