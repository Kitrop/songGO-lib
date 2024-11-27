package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Подгрузка env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	
}