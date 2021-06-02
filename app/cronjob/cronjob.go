package cronjob

import (
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
	testArray := strings.Split(input, "Name:")
	val := len(testArray)
	snapshotLimit := os.Getenv("SnapshotLimit")
	limit, _ := strconv.Atoi(snapshotLimit)
	for i := 1; i < val-limit; i++ {
		existingSnapshotName := strings.Split(testArray[i], " ")[1]
		_, err := ExecCommandOnHost("VBoxManage snapshot " + vmname + " delete " + existingSnapshotName)
		if err != nil {
			return err
		}
	}
	return nil
}
