package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) routes() {
	s.router.With(devOnly).Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("docs/"))))
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Use(accessControl)
		r.Use(validateToken(s.oauth2Server))
		r.Route("/users", func(r chi.Router) {
			r.Get("/", s.userHandler.HandleGetUsers())
			r.Get("/{email}", s.userHandler.HandleGetUserByEmail())
			r.Post("/", s.userHandler.HandleCreateUser())
		})
	})
}
