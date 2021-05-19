package main

import (
	"fmt"
	"sync"

	"app/CRON"
	"app/DB"
	"app/Routes"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load(".env")
	wg := new(sync.WaitGroup)
	wg.Add(3)

	r := Routes.SetupRouter()
	DB.EstablishConnection()

	go func() {
		defer wg.Done()
		CRON.C = cron.New()
		CRON.C.Start()
	}()

	go func() {
		defer wg.Done()
		CRON.Jobs()
	}()

	go func() {
		defer wg.Done()
		r.Run(":5001")
		fmt.Println("Listen and Server in 0.0.0.0:5001")
	}()

	wg.Wait()
}
