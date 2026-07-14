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

	repo, err := repository.AddClient(postgresCreds)
	if err != nil {
		log.Fatalln(ewrap.Wrap("Couldn't create repo client", err))
	}

	authService := service.NewAuthService(repo, []byte(os.Getenv("JWT_SECRET")))

	transactionsService := service.NewTransactionService(repo)

	categoryService, err := service.NewCategoryService(repo)
	if err != nil {
		log.Fatalln(ewrap.Wrap("Couldn't get system categories", err))
	}

	budgetService := service.NewBudgetService(repo)

	server := server.AddServer(os.Getenv("HOST"), authService, transactionsService, categoryService, budgetService)

	log.Println("Backend is running")

	server.Start(os.Getenv("BACKEND_PORT"))
}
