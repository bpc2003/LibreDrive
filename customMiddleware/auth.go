package customMiddleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		auth = strings.Split(auth, " ")[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(auth, claims, func(tok *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			initContext := r.Context()
			for key, val := range claims {
				initContext = context.WithValue(initContext, key, val)
			}
			next.ServeHTTP(w, r.WithContext(initContext))
		}
	})
}
