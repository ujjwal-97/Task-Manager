package cronjob

import (
	"log"
	"os"

	"gopkg.in/robfig/cron.v2"
)

func CreateSystemBackupCron(c *cron.Cron) error {
	_, err := c.AddFunc("@midnight", func() {
		out, err := CreateSystemBackup()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(out))
	})
	return err
}

func CreateSystemBackup() (string, error) {
	cmd := "echo " + os.Getenv("hostpassword") + "|" + " sudo -S timeshift --create"
	out, err := ExecCommandOnHost(cmd)
	return out, err
}
