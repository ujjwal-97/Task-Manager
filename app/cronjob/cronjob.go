package cronjob

import (
	"app/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var C *cron.Cron

func TakeSnapshot(vmname string, snapshotName string) (string, error) {

	datetime := time.Now().Format(time.RFC3339)
	snapshotName += "_" + datetime

	output, _ := ExecCommandOnHost("VBoxManage snapshot " + vmname + " list")

	if err := checkSnapshotLimit(output, vmname); err != nil {
		return "", err
	}

	output, err := ExecCommandOnHost("VBoxManage snapshot " + vmname + " take " + snapshotName)
	if err != nil {
		return "Error Executing Command", err
	}

	return string(output), nil
}

func checkSnapshotLimit(input string, vmname string) error {
	snapDetails := strings.Split(input, "Name:")
	val := len(snapDetails)
	snapshotLimit := os.Getenv("SnapshotLimit")
	limit, _ := strconv.Atoi(snapshotLimit)
	for i := 1; i < val-limit; i++ {
		existingSnapshotName := strings.Split(snapDetails[i], " ")[1]
		_, err := ExecCommandOnHost("VBoxManage snapshot " + vmname + " delete " + existingSnapshotName)
		if err != nil {
			return err
		}
	}
	return nil
}

func ScheduleSnapshot(schedule *models.Schedule, vmname string, snapshotName string) (string, error) {
	timming := formatTime(schedule.Hour) + ":" + formatTime(schedule.Minute) + " "
	if schedule.Month != 0 {
		timming += " " + formatTime(schedule.Month) + formatTime(schedule.Day) + formatTime(time.Now().Year()%100)
	}
	timestamp := time.Now().Format(time.RFC3339)
	snapshotName += timestamp
	cmd := "echo vboxmanage snapshot " + vmname + " take " + snapshotName + " | at " + timming
	output, err := ExecCommandOnHost(cmd)
	if err != nil {
		return "Error Executing Command", err
	}

	return string(output), nil
}
func formatTime(value int) string {
	if value < 10 {
		return "0" + strconv.Itoa(value)
	}
	return strconv.Itoa(value)
}

func CreateCronExpression(s *models.ScheduleSnapshot) string {
	ans := ""
	ans += createStringCron(s.Schedule.Minute, 60, s.Interval.Minute) + " "
	ans += createStringCron(s.Schedule.Hour, 24, s.Interval.Hour) + " "
	ans += createStringCron(s.Schedule.Day, 31, s.Interval.Day) + " "
	ans += createStringCron(s.Schedule.Month, 12, s.Interval.Month) + " "
	weekday := ""
	if s.Schedule.Weekday >= 0 && s.Schedule.Weekday <= 7 {
		weekday += strconv.Itoa(s.Schedule.Weekday)
	} else {
		weekday += "*"
	}
	ans += weekday

	return ans
}

func createStringCron(value int, mod int, factor int) string {
	newTimeUnit := ""
	if value > 0 {
		newTimeUnit += strconv.Itoa(value % mod)
	} else {
		newTimeUnit += "*"
	}
	if factor > 1 {
		newTimeUnit += "/" + strconv.Itoa(factor)
	}
	return newTimeUnit
}
