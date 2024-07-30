package customMiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"libredrive/types"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusForbidden)
			enc.Encode(types.ErrStruct{Success: false, Msg: "Forbidden"})
			return
		}
		auth = strings.Split(auth, " ")[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(auth, claims, func(tok *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
		} else {
			initContext := r.Context()
			for key, val := range claims {
				initContext = context.WithValue(initContext, key, val)
			}
			next.ServeHTTP(w, r.WithContext(initContext))
		}
	})
}
