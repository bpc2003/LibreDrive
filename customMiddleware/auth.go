package customMiddleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, _ := r.Cookie("auth")
		if auth == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		attrs := strings.Split(auth.Value, "&")
		id, _ := strconv.Atoi(attrs[0])
		isAdmin, _ := strconv.ParseBool(attrs[1])

		initContext := context.WithValue(r.Context(), "id", id)
		initContext = context.WithValue(initContext, "isAdmin", isAdmin)
		next.ServeHTTP(w, r.WithContext(initContext))
	})
}
