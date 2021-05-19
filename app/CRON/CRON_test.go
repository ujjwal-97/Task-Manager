package CRON

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gopkg.in/robfig/cron.v2"
)

func TestCheckUpdateJob(t *testing.T) {
	c := cron.New()
	c.Start()

	err := checkUpdateJob(c)
	assert.NoError(t, err)

	os.Setenv("Password", "-")
	err = checkUpdateJob(c)
	assert.NoError(t, err)
}

func TestJobs(t *testing.T) {
	err := Jobs()
	assert.NoError(t, err)
}

func TestHealthCheckJob(t *testing.T) {
	c := cron.New()
	c.Start()

	err := healthCheckJob(c)
	assert.NoError(t, err)
}

func TestCheckSystemHealth(t *testing.T) {
	godotenv.Load("../.env")

	_, err := CheckSystemHealth()
	assert.Error(t, err)

	os.Setenv("healthcheckScript", "../healthCheck.sh")
	_, err = CheckSystemHealth()
	assert.NoError(t, err)
}

func TestTakeSnapshot(t *testing.T) {
	godotenv.Load("../.env")

	hostIP := os.Getenv("hostip")
	hostUser := os.Getenv("hostusername")
	hostPassword := os.Getenv("hostpassword")

	conn, err := Connect(hostIP, hostUser, hostPassword)
	assert.NoError(t, err)

	_, err = conn.SendCommands("ls")
	assert.NoError(t, err)

	_, err = TakeSnapshot("ubuntu2", "ubuntu2")
	assert.NoError(t, err)

	_, err = TakeSnapshot("nil", "nil")
	assert.Error(t, err)
}
