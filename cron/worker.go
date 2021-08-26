package cron

import (
	"fmt"
	"fssync/config"
	"fssync/core"
	"fssync/util"
	"strings"
)

func Worker() {
	var src []string
	diffFiles, err := util.GetDiffFile()
	if err != nil {
		config.Logger.Error("Get diff file error:", err.Error())
	} else if diffFiles != nil {
		config.Logger.Info("get files:", diffFiles)
	}
	// 按文件路径信息执行同步
	for _, file := range diffFiles {
		config.Logger.Infof("Start sync file:%s", file)
		dirs := strings.Split(file, "/")
		if len(dirs) > 3 {
			mod, env, appName := dirs[0], dirs[1], dirs[2]
			srcPath := fmt.Sprintf("%s/%s/%s/%s/", config.Config.RepoDir, mod, env, appName)

			if !util.Contains(&src, &srcPath) {
				src = append(src, srcPath)
				ipList, err := util.Getip(&env, &appName)
				if err != nil {
					config.Logger.Error("Get ip error:", err.Error())
					return
				}
				config.Logger.Infof("[Sync info]src:%s;mod:%s;env:%s;app:%s;iplist:%s", srcPath, mod, env, appName, ipList)
				for _, ip := range ipList {
					go core.SyncCron(srcPath, ip)
				}
			}
		}
	}
}
