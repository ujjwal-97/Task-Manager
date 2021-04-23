package main

import (
	"fmt"
	"log"

	"./Controllers/Connect"
	"./Routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := Routes.SetupRouter()

	Connect.EstablishConnection()

	r.Run(":5001")

	fmt.Println("Listen and Server in 0.0.0.0:5001")
}
