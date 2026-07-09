package main

import (
	"log"
	"os"

	ewrap "finstat/internal/lib"
	"finstat/internal/repository"
	"finstat/internal/server"
	"finstat/internal/service"
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

	repo, err := repository.NewClient(postgresCreds)
	if err != nil {
		log.Fatalln(ewrap.Wrap("Couldn't create repo client", err))
	}

	authService := service.NewAuthService(repo, []byte(os.Getenv("JWT_SECRET")))
	transactionsService := service.NewTransactionService(repo)

	server := server.NewServer(os.Getenv("HOST"), authService, transactionsService)

	log.Println("Backend is running")

	server.Start(os.Getenv("BACKEND_PORT"))
}
