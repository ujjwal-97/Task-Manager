package cronjob

import (
	"time"

	"github.com/robfig/cron/v3"
)

var C *cron.Cron

func TakeSnapshot(vmname string, snapshotName string) (string, error) {

	datetime := time.Now().Format(time.RFC3339)
	snapshotName += "_" + datetime

	output, err := ExecCommandOnHost("VBoxManage snapshot " + vmname + " take " + snapshotName)
	if err != nil {
		return "Error Executing Command", err
	}

	return string(output), nil
}
