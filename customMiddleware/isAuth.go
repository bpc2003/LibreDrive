package customMiddleware

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")
		if userId == "" && r.Context().Value("isAdmin").(bool) {
			next.ServeHTTP(w, r)
		} else if userId != "" {
			userId_int, err := strconv.Atoi(userId)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			if userId_int == r.Context().Value("id").(int) || r.Context().Value("isAdmin").(bool) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
