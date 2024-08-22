package controllers

import (
	"context"
	"net/http"

	"libredrive/templates"
)

func Nav(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	isAdmin := r.Context().Value("isAdmin").(bool)

	nav := templates.Nav(id, isAdmin)
	nav.Render(context.Background(), w)
}
