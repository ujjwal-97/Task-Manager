package main

import (
	"fmt"
	"sync"

	"app/cronjob"
	"app/db"
	"app/routes"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load(".env")
	wg := new(sync.WaitGroup)
	wg.Add(3)

	r := routes.SetupRouter()
	db.EstablishConnection()

	go func() {
		defer wg.Done()
		cronjob.C = cron.New()
		cronjob.C.Start()
	}()

	go func() {
		defer wg.Done()
		cronjob.Jobs()
	}()

	go func() {
		defer wg.Done()
		r.Run(":5001")
		fmt.Println("Listen and Server in 0.0.0.0:5001")
	}()

	wg.Wait()
}
