package customMiddleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"libredrive/types"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		userId := chi.URLParam(r, "userId")
		if userId == "" && r.Context().Value("isAdmin").(bool) {
			next.ServeHTTP(w, r)
		} else if userId != "" {
			userId_int, err := strconv.Atoi(userId)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				enc.Encode(types.ErrStruct{Success: false, Msg: "Invalid ID"})
				return
			}

			if float64(userId_int) == r.Context().Value("id").(float64) || r.Context().Value("isAdmin").(bool) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				enc.Encode(types.ErrStruct{Success: false, Msg: "Forbidden"})
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			enc.Encode(types.ErrStruct{Success: false, Msg: "Forbidden"})
		}
	})
}
