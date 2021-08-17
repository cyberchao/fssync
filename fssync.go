package main

import (
	"fssync/api"
	"fssync/config"
	mycron "fssync/cron"
	"fssync/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	config.VP = config.Viper("./config.yaml")
	logger.InitLogger()
	// mycron.Worker()
	mycron.Cron()
	router := gin.Default()
	router.GET("/sync", api.SyncFunc)
	router.Run(":8080")
}
