package main

import (
	"fssync/api"
	"fssync/config"
	mycron "fssync/cron"
	"fssync/logger"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	config.VP = config.Viper("./config.yaml")
	logger.InitLogger()
	// mycron.Worker()
	mycron.Cron()
	router := gin.Default()
	// k8s通过此静态文件下载服务下载配置文件
	router.Use(static.Serve("/", static.LocalFile(config.Config.RepoDir+"/env", true)))
	// 获取文件列表
	router.GET("/getfile", api.GetFileFunc)
	// 手动同步
	router.GET("/sync", api.SyncFunc)
	// 配置修改
	router.POST("/modify", api.ModifyFunc)
	router.Run(":8080")
}
