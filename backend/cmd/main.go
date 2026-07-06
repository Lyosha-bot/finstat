package main

import (
	"log"
	"os"

	ewrap "auth.my-financials/internal/lib"
	"auth.my-financials/internal/repository"
	"auth.my-financials/internal/server"
)

func main() {
	log.Println("Starting backend")

	postgresCreds := repository.Credentials{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DB_name:  os.Getenv("DB_NAME"),
	}

	server, err := server.NewServer(postgresCreds, os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalln(ewrap.Wrap("Couldn't start backend", err))
	}

	log.Println("Backend is running")

	server.Start(os.Getenv("BACKEND_PORT"))
}
