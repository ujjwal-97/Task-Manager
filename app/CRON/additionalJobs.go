package CRON

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/robfig/cron.v2"
)

func Jobs() error {

	c := cron.New()
	c.Start()
	err := checkUpdateJob(c)
	if err != nil {
		log.Println(err.Error())
	}
	err = healthCheckJob(c)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func checkUpdateJob(c *cron.Cron) error {
	_, err := c.AddFunc("@midnight", func() {
		password := os.Getenv("Password")
		cmd := "sudo -S <<<" + password + " apt-get update && sudo apt upgrade"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(out))
	})
	return err
}

func healthCheckJob(c *cron.Cron) error {
	schedule := "0 0 0 1 * ?"
	//schedule := "0 * * * * *"
	_, err := c.AddFunc(schedule, func() {
		status, _ := CheckSystemHealth()
		log.Println(status)
	})
	return err
}

func CheckSystemHealth() (string, error) {
	filename := os.Getenv("healthcheckScript")
	exec.Command("chmod", "+x", filename)
	cmd := "./" + filename
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", err.Error()), err
	}
	return string(out), nil
}
