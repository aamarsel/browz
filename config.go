package main

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Файл .env не найден, используем системные переменные")
	}
}
