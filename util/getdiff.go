package util

import (
	"bytes"
	"fssync/config"
	"os"
	"os/exec"
	"strings"
)

// 获取与上个版本有差异的文件列表
func GetDiffFile() ([]string, error) {
	os.Chdir(config.RepoDir)
	cmd := exec.Command("git", "pull", "origin", "main")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		config.Logger.Errorf("git pull error:%s:%s", err, stderr.String())
		return nil, err
	} else {
		config.Logger.Info("git pull: " + strings.Trim(out.String(), "\n"))
		if !strings.Contains(out.String(), "Already up to date") {
			// git log -p -1 --oneline
			out, _ := exec.Command("git", "diff", "head^", "--name-only").Output()
			files := strings.Split(strings.Trim(string(out), "\n"), "\n")
			return files, nil
		} else {
			config.Logger.Infof("No diff file")
			return nil, nil
		}
	}
}
