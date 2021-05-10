package CRON

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gopkg.in/robfig/cron.v2"
)

func TestCheckUpdateJob(t *testing.T) {
	c := cron.New()
	c.Start()
	err := checkUpdateJob(c)
	if err != nil {
		t.Errorf(err.Error())
	}
	os.Setenv("Password", "12345678")
	err = checkUpdateJob(c)
	if err != nil {
		t.Errorf(err.Error())
	}
}
func TestJobs(t *testing.T) {
	if err := Jobs(); err != nil {
		t.Errorf(err.Error())
	}
}
func TestHealthCheckJob(t *testing.T) {
	c := cron.New()
	c.Start()
	err := healthCheckJob(c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCheckSystemHealth(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	status, err := CheckSystemHealth()
	if err == nil {
		t.Errorf("%v", status)
	}
	os.Setenv("healthcheckScript", "../healthCheck.sh")
	status, err = CheckSystemHealth()
	if err != nil {
		t.Errorf("%v", status)
	}
}
