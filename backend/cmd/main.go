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

	repos, err := repository.New(postgresCreds)
	if err != nil {
		log.Fatalln(lib.Ewrap("Couldn't create repo client", err))
	}

	services := service.New(repos, []byte(os.Getenv("JWT_ACCESS_SECRET")), []byte(os.Getenv("JWT_REFRESH_SECRET")))

	server := server.NewServer(os.Getenv("HOST"), services)

	log.Println("Backend is running")

	server.Start(os.Getenv("BACKEND_PORT"))
}
