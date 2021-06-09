package cronjob

import (
	"app/cronjob"
	"app/models"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gopkg.in/robfig/cron.v2"
)

func TestCheckUpdateJob(t *testing.T) {
	godotenv.Load("../../.env")
	c := cron.New()
	c.Start()

	err := cronjob.CheckUpdateJob(c)
	assert.NoError(t, err)

	os.Setenv("Password", "-")
	err = cronjob.CheckUpdateJob(c)
	assert.NoError(t, err)
}
func TestExecCommandOnHost(t *testing.T) {
	cmd := "hostname"
	res, err := cronjob.ExecCommandOnHost(cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestJobs(t *testing.T) {
	err := cronjob.Jobs()
	assert.NoError(t, err)
}

func TestHealthCheckJob(t *testing.T) {
	c := cron.New()
	c.Start()

	err := cronjob.HealthCheckJob(c)
	assert.NoError(t, err)
}

func TestCheckSystemHealth(t *testing.T) {
	godotenv.Load("../../.env")

	_, err := cronjob.CheckSystemHealth()
	assert.Error(t, err)

	os.Setenv("healthcheckScript", "../../healthCheck.sh")
	_, err = cronjob.CheckSystemHealth()
	assert.NoError(t, err)
}

func TestTakeSnapshot(t *testing.T) {
	godotenv.Load("../../.env")

	hostIP := os.Getenv("hostip")
	hostUser := os.Getenv("hostusername")
	hostPassword := os.Getenv("hostpassword")

	conn, err := cronjob.Connect(hostIP, hostUser, hostPassword)
	assert.NoError(t, err)

	_, err = conn.SendCommands("ls")
	assert.NoError(t, err)

	_, err = cronjob.TakeSnapshot("ubuntu2", "ubuntu2")
	assert.NoError(t, err)

	_, err = cronjob.TakeSnapshot("nil", "nil")
	assert.Error(t, err)
}

func TestScheduleSnapshot(t *testing.T) {
	godotenv.Load("../../.env")

	hostIP := os.Getenv("hostip")
	hostUser := os.Getenv("hostusername")
	hostPassword := os.Getenv("hostpassword")

	conn, err := cronjob.Connect(hostIP, hostUser, hostPassword)
	assert.NoError(t, err)

	_, err = conn.SendCommands("ls")
	assert.NoError(t, err)

	timming := models.Schedule{
		Minute:  time.Now().Minute() + 1,
		Hour:    time.Now().Hour(),
		Day:     0,
		Month:   0,
		Weekday: -1,
	}
	_, err = cronjob.ScheduleSnapshot(&timming, "ubuntu2", "ubuntu2")
	assert.NoError(t, err)
	timming = models.Schedule{
		Minute:  time.Now().Minute() + 2,
		Hour:    time.Now().Hour(),
		Day:     time.Now().Day(),
		Month:   timming.Month,
		Weekday: -1,
	}
	_, err = cronjob.ScheduleSnapshot(&timming, "ubuntu2", "ubuntu2")
	assert.NoError(t, err)
}

func TestCreateCronExpression(t *testing.T) {
	s := models.ScheduleSnapshot{
		Periodic: false,
		Schedule: &models.Schedule{
			Minute:  0,
			Hour:    0,
			Day:     0,
			Month:   0,
			Weekday: -1,
		},
		Interval: &models.Interval{
			Minute: 0,
			Hour:   0,
			Day:    0,
			Month:  0,
		},
	}
	res := cronjob.CreateCronExpression(&s)
	ans := "* * * * *"
	assert.Equal(t, ans, res)

}
