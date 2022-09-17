package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"pet-store/api"
	"pet-store/database"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Unable to load environment variables file")
	}
	database.Init()
	db := database.DBCon
	defer db.Close()

	http.ListenAndServe(":8080", api.GetRouter())
}
