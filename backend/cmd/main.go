package main

import (
	"log"
	"os"

	"finstat/internal/lib"
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

	repo, err := repository.InsertClient(postgresCreds)
	if err != nil {
		log.Fatalln(lib.Ewrap("Couldn't create repo client", err))
	}

	authService := service.NewAuthService(repo, []byte(os.Getenv("JWT_ACCESS_SECRET")), []byte(os.Getenv("JWT_REFRESH_SECRET")))

	transactionsService := service.NewTransactionService(repo)

	categoryService := service.NewCategoryService(repo)

	budgetService := service.NewBudgetService(repo)

	server := server.InsertServer(os.Getenv("HOST"), authService, transactionsService, categoryService, budgetService)

	log.Println("Backend is running")

	server.Start(os.Getenv("BACKEND_PORT"))
}
