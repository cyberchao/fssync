package core

import (
	"fssync/config"
	"os/exec"
)

func SyncCron(srcPath, ip *string) {
	err := exec.Command("/usr/bin/rsync", "-avz", "--timeout=3", "--owner=wls81", "--group=wls", *srcPath, "root@"+*ip+":/").Run()
	if err != nil {
		config.Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", *srcPath, *ip, err.Error())
	} else {
		config.Logger.Infof("Rsync success:rsync -az %s* %s:/", *srcPath, *ip)
	}
}

func SyncHttp(srcPath, ip *string, ch chan string) {
	out, err := exec.Command("/usr/bin/rsync", "-avz", "--timeout=3", "--owner=wls81", "--group=wls", *srcPath, "root@"+*ip+":/").Output()
	if err != nil {
		config.Logger.Infof("Rsync error:%s", out)
		config.Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", *srcPath, *ip, err.Error())
		ch <- *ip + ":" + string(out)
	} else {
		config.Logger.Infof("Rsync success:rsync -az %s* %s:/", *srcPath, *ip)
		ch <- *ip + ":success"
	}
}
