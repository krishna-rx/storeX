package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"main.go/database"
	"main.go/server"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading from .env file")
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	ConStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	database.InitDBAndMigrate(ConStr)
	defer database.CloseDB()
	s := server.SetUpRoutes(":8080")
	s.Start()
}
