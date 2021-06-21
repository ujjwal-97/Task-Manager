package cronjob

import (
	"log"
	"os"

	"gopkg.in/robfig/cron.v2"
)

func CreateDumpCron(c *cron.Cron) error {
	_, err := c.AddFunc("@midnight", func() {
		out, err := CreateDump()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(out))
	})
	return err
}

func CreateDump() (string, error) {

	cmd := "mongodump --host " + os.Getenv("DB_HOST") + " --port " + os.Getenv("DB_PORT") + " --authenticationDatabase admin " + " -u " + os.Getenv("DB_USERNAME") + " -p " + os.Getenv("DB_PASSWORD") + " -d " + os.Getenv("DB_NAME") + " --archive=" + os.Getenv("BackupPath") + " --gzip"
	out, err := ExecCommandOnHost(cmd)
	return out, err
}

func RestoreDump() (string, error) {

	cmd := "mongorestore --host " + os.Getenv("DB_HOST") + " --port " + os.Getenv("DB_PORT") + " --authenticationDatabase admin " + " -u " + os.Getenv("DB_USERNAME") + " -p " + os.Getenv("DB_PASSWORD") + " -d " + os.Getenv("DB_NAME") + " -c task" + "  --gzip --archive=" + os.Getenv("BackupPath")
	out, err := ExecCommandOnHost(cmd)
	return out, err
}
