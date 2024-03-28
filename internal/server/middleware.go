package server

import (
	"goddd/internal"
	"goddd/internal/oauth2"
	"goddd/pkg/errx"
	"net/http"
	"os"
)

var ErrUnauthorized = errx.NewErrorf(internal.CodeInvalidToken, "invalid access token")

func devOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") != "dev" {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte("route does not exist")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func accessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(oauth2Server *oauth2.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := oauth2Server.ValidateToken(r); err != nil {
				encodeError(w, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
