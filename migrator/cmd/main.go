package main

import (
	"log"
	"os"
	"postgres-migrator/internal/migrator"
	"postgres-migrator/internal/repository"
)

func main() {
	creds := repository.Credentials{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DB_name:  os.Getenv("DB_NAME"),
	}

	log.Println("Migrations start")
	if err := migrator.ApplyMigrations(creds); err != nil {
		log.Println("Migrations failed")
		log.Fatalln(err)
	}
	log.Println("Migrations complete")
}
