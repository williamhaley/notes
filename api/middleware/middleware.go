package middleware

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt"
)

func New(serverSecret string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(serverSecret), nil
			})
			if err != nil {
				log.Error("error parsing token", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("bad request"))
				return
			}

			if err, ok := token.Claims.(jwt.MapClaims); !ok {
				log.Error("invalid claims", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("bad request"))
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
