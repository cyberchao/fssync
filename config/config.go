package config

import "go.uber.org/zap"

const (
	Workdir  = "/Users/pangru/Documents/golang/fssync/"
	Logfile  = "/Users/pangru/Documents/golang/fssync/syncer.log"
	RepoDir  = "/Users/pangru/Documents/golang/fssync/serverfiles" // git本地仓库目录
	Interval = "10"                                                //cron 间隔
	Timeout  = "3"
	Owner    = "wls81"
	Group    = "wls"
)

var (
	Logger *zap.SugaredLogger
)
