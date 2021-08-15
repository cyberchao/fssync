package cron

import (
	"fmt"
	"fssync/config"

	"github.com/robfig/cron"
)

func Cron() {
	cron := cron.New()

	cron.AddFunc(fmt.Sprintf("@every %ss", config.Interval), func() {
		config.Logger.Infof("start cron")
		Worker()
		config.Logger.Infof("end cron")
	})

	cron.Start()
}
