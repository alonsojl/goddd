package main

import (
	"fmt"

	"goddd/internal/domain/user"
	"goddd/internal/mysql"
	"goddd/internal/oauth2"
	"goddd/internal/server"
	"goddd/pkg/logger"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

// @title Go RESTful API
// @version 1.0

// @schemes 		http https
// @basePath 		/api/v1
// @description 	Testing Swagger APIs.
// @termsOfService 	http://swagger.io/terms/

// @contact.name 	API Support
// @contact.url		http://www.swagger.io/support
// @contact.email 	support@swagger.io

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @securityDefinitions.apiKey JWT
// @in	 header
// @name Authorization
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		host     = os.Getenv("API_HOST")
		port     = os.Getenv("API_PORT")
		httpAddr = fmt.Sprintf("%s:%s", host, port)
	)

	db, err := mysql.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	logger, err := logger.New()
	if err != nil {
		return err
	}

	userRepository := mysql.NewUserRepository(logger, db)
	userService := user.NewUserService(logger, userRepository)
	userHandler := server.NewUserHandler(logger, userService)
	oauth2Server := oauth2.NewServer(logger, db)
	router := chi.NewRouter()
	srv := server.New(httpAddr, router, logger, oauth2Server, userHandler)

	if err = <-srv.Run(); err != nil {
		return err
	}

	logger.Info("terminated")
	return nil
}
