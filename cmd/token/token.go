package main

import (
	"encoding/json"
	"fmt"
	"goddd/internal/mysql"
	"goddd/internal/oauth2"

	_ "github.com/go-sql-driver/mysql"

	"goddd/pkg/logger"
	"log"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		log.Printf("error initializing logger: %v\n", err)
		return
	}

	db, err := mysql.Connection()
	if err != nil {
		log.Printf("error connecting to mysql database: %v\n", err)
		return
	}

	server := oauth2.NewServer(logger, db)
	credentials, err := server.CreateCredentials()
	if err != nil {
		log.Printf("error creating credentials: %v\n", err)
		return
	}

	token, err := server.CreateAccessToken(credentials)
	if err != nil {
		log.Printf("error creating token: %v\n", err)
		return
	}

	data := map[string]string{
		"token": token,
	}
	result, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("error converting to JSON:", err)
		return
	}

	fmt.Println(string(result))
}
