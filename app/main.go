package main

import (
	"fmt"
	"log"
	"sync"

	"app/CRON"
	"app/DB"
	"app/Routes"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)

	r := Routes.SetupRouter()
	DB.EstablishConnection()

	go func() {
		CRON.C = cron.New()
		CRON.C.Start()
		wg.Done()
	}()

	go func() {
		r.Run(":5001")
		fmt.Println("Listen and Server in 0.0.0.0:5001")
		wg.Done()
	}()

	wg.Wait()
}
