package cron

import (
	"fmt"
	"fssync/config"

	"github.com/robfig/cron"
)

func Cron() {
	cron := cron.New()

	cron.AddFunc(fmt.Sprintf("@every %ss", config.Config.Interval), func() {
		Worker()
	})

	cron.Start()
}
