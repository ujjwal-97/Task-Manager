package CRON

import (
	"testing"

	"github.com/robfig/cron/v3"
)

func TestCron(t *testing.T) {
	C = cron.New()
	if C == nil {
		t.Error()
	}
}
