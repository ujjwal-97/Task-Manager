package CRON

import (
	"testing"

	"github.com/robfig/cron"
)

func TestCron(t *testing.T) {
	C = cron.New()
	if C == nil {
		t.Error()
	}
}
