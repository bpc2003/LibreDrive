// middleware - middleware for handling private routes
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Auth - makes sure a user is logged in.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, _ := r.Cookie("auth")
		if auth == nil {
			w.Header().Set("HX-Redirect", "/login.html")
			return
		}
		attrs := strings.Split(auth.Value, "&")
		id, _ := strconv.Atoi(attrs[0])
		isAdmin, _ := strconv.ParseBool(attrs[1])

		initContext := context.WithValue(r.Context(), "id", id)
		initContext = context.WithValue(initContext, "isAdmin", isAdmin)
		initContext = context.WithValue(initContext, "key", attrs[2])
		next.ServeHTTP(w, r.WithContext(initContext))
	})
}
