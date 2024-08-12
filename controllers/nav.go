package controllers

import (
	"context"
	"net/http"

	"libredrive/templates"
)

func Nav(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value("isAdmin").(bool)

	nav := templates.Nav(isAdmin)
	nav.Render(context.Background(), w)
}
