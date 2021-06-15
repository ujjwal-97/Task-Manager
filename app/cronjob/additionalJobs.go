package cronjob

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
	err := CheckUpdateJob(c)
	if err != nil {
		log.Println(err.Error())
	}
	err = HealthCheckJob(c)
	if err != nil {
		log.Println(err.Error())
	}
	err = StartUpPreviousCron(c)
	if err != nil {
		log.Println(err.Error())
	}
	return err

}

func CheckUpdateJob(c *cron.Cron) error {
	_, err := c.AddFunc("@midnight", func() {
		cmd := "apt update && apt upgrade"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(out))
	})
	return err
}

func HealthCheckJob(c *cron.Cron) error {
	schedule := "0 0 0 1 * ?"
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
