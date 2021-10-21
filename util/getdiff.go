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
	os.Chdir(config.Config.RepoDir)
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
		// 判断是否有更新
		if !strings.Contains(out.String(), "Already up to date") {
			config.Logger.Info("git pull: " + strings.Trim(out.String(), "\n"))
			// git log -p -1 --oneline 获取最近一次更新的详细内容变化
			out, _ := exec.Command("git", "diff", "HEAD^", "--name-only").Output()
			files := strings.Split(strings.Trim(string(out), "\n"), "\n")
			return files, nil
		} else {
			config.Logger.Infof("No changed files")
			return nil, nil
		}
	}
}
