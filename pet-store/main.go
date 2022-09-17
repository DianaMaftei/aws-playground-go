package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"pet-store/api"
	"pet-store/database"
)

func main() {
	env := os.Getenv("ENV")
	configFile := ".env"
	if env == "LOCAL" {
		configFile = "local.env"
	}

	err := godotenv.Load(configFile)
	if err != nil {
		fmt.Println("Unable to load .env file")
		os.Exit(1)
	}
	database.Init()
	db := database.DBCon
	defer db.Close()

	http.ListenAndServe(":8080", api.GetRouter())
}
