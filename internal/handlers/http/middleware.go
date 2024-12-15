package http

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
)

func validateToken(signature string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerToken := r.Header.Get("Token")
			if headerToken == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			token, err := jwt.Parse(headerToken, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return false, errors.New("invalid Token")
				}
				return []byte(signature), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !token.Valid {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
