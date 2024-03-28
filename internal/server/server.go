package server

import (
	"context"
	"errors"
	"goddd/internal/oauth2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpAddr     string
	router       *chi.Mux
	logger       *logrus.Logger
	oauth2Server *oauth2.Server
	userHandler  UserHandler
}

func New(httpAddr string, router *chi.Mux, logger *logrus.Logger, oauth2Server *oauth2.Server, userHandler UserHandler) *Server {
	s := &Server{
		httpAddr:     httpAddr,
		router:       router,
		logger:       logger,
		oauth2Server: oauth2Server,
		userHandler:  userHandler,
	}
	s.routes()
	return s
}

func (s *Server) Run() <-chan error {
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.router,
	}

	notify := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		s.logger.Info("shutdown signal received")

		timeout := 10 * time.Second
		ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
		defer func() {
			stop()
			cancel()
			close(notify)
		}()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctxTimeout); err != nil {
			notify <- err
		}
		s.logger.Info("shutdown completed")
	}()

	go func() {
		s.logger.Infof("listening and serving %s", s.httpAddr)
		// ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			notify <- err
		}
	}()

	return notify
}
