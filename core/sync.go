package core

import (
	"fssync/config"
	"os/exec"
)

func SyncCron(srcPath, ip string) {
	err := exec.Command("/usr/bin/rsync", "-avz", "--timeout="+config.Config.Timeout, "--owner="+config.Config.Owner, "--group="+config.Config.Group, srcPath, ip+":/").Run()
	if err != nil {
		config.Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", srcPath, ip, err.Error())
	} else {
		config.Logger.Infof("Rsync success:rsync -az %s* %s:/", srcPath, ip)
	}
}

func SyncHttp(srcPath, ip string, ch chan string) {
	out, err := exec.Command("/usr/bin/rsync", "-avz", "--timeout="+config.Config.Timeout, "--owner="+config.Config.Owner, "--group="+config.Config.Group, srcPath, ip+":/").Output()
	if err != nil {
		config.Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", srcPath, ip, err.Error())
		ch <- ip + ":" + string(out)
	} else {
		config.Logger.Infof("Rsync success:rsync -az %s* %s:/", srcPath, ip)
		ch <- ip + ":success"
	}
}
