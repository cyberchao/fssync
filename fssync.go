package main

import (
	"fssync/api"
	mycron "fssync/cron"
	"fssync/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()
	// mycron.Worker()
	mycron.Cron()
	router := gin.Default()
	router.GET("/sync", api.SyncFunc)
	router.Run(":8080")
}
