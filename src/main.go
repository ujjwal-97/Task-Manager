package main

import (
	"fmt"

	"./Controllers/Connect"
	"./Routes"
)

func main() {
	r := Routes.SetupRouter()
	Connect.EstablishConnection()
	r.Run(":8080")
	fmt.Println("Listen and Server in 0.0.0.0:8080")
}
