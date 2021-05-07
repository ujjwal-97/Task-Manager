package CRON

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/robfig/cron.v2"
)

func Jobs() {

	c := cron.New()
	c.Start()

	c.AddFunc("@midnight", func() {
		password := os.Getenv("Password")
		cmd := "sudo -S <<<" + password + " apt-get update && sudo apt upgrade"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(out))
	})

	schedule := "0 21 15 1 * ?"
	//schedule := "0 * * * * *"
	c.AddFunc(schedule, func() {
		log.Println(CheckSystemHealth())
	})
}

func CheckSystemHealth() string {
	filename := os.Getenv("healthcheckScript")
	exec.Command("chmod", "+x", filename)
	cmd := "./" + filename
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", err.Error())
	}
	return string(out)
}
