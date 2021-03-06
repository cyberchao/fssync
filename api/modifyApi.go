package api

import (
	"bytes"
	"errors"
	"fmt"
	"fssync/config"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// 发布平台可以通过此接口修改配置文件
func ModifyFunc(c *gin.Context) {
	var requestData Request
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "msg": "data struct error:" + err.Error()})
		return
	}
	msg, err := Modify(requestData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "msg": "edit error:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "true",
		"msg":    msg,
	})
}

type Request struct {
	AppName  string            `json:"appName"`
	EnvName  string            `json:"envName"`
	Path     string            `json:"path"`
	Filename string            `json:"filename"`
	Operate  string            `json:"operate"`
	Datas    map[string]string `json:"datas"`
}

// type Request struct {
// 	Msg []ModifyData
// }

func Modify(requestData Request) (string, error) {
	os.Chdir(config.Config.Easypub)
	cmd := exec.Command("git", "pull", "origin", "main")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		config.Logger.Errorf("git pull error:%s:%s", err, stderr.String())
		return "", err
	}

	filePath := fmt.Sprintf("%s/env/%s/%s/%s/%s", config.Config.Easypub, requestData.EnvName, requestData.AppName, requestData.Path, requestData.Filename)
	fi, err := os.Lstat(filePath)
	if err != nil {
		return "", err
	}

	perm := fi.Mode().Perm()
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(input), "\n")

	// 增加行
	if requestData.Operate == "add" {
		for k, v := range requestData.Datas {
			lineArray := []string{k, v}
			newLine := strings.Join(lineArray, "=")
			lines = append(lines, newLine)
		}
		// 删除或修改行
	} else {
		for i, line := range lines {
			lineArray := strings.SplitN(line, "=", 2)
			for k, v := range requestData.Datas {
				if strings.Contains(line, k) {
					if requestData.Operate == "edit" {
						lineArray[1] = v
						newLine := strings.Join(lineArray, "=")
						lines[i] = newLine
					} else if requestData.Operate == "del" {
						lines[i] = ""
					} else {
						return "", errors.New("operation error")
					}
				}
			}
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), perm)
	if err != nil {
		exec.Command("git", "reset", "--hard").Run()
		return "", err
	}
	os.Chdir(config.Config.Easypub)
	cmd = exec.Command("git", "add", "./*")
	cmd.Stdout, cmd.Stderr = &out, &stderr
	err = cmd.Run()
	if err != nil {
		config.Logger.Errorf("git add error:%s:%s:%s", err, stderr.String(), out.String())
		exec.Command("git", "reset", "--hard").Run()
		return "", err
	}
	cmd = exec.Command("git", "commit", "-m", " easypub")
	cmd.Stdout, cmd.Stderr = &out, &stderr
	err = cmd.Run()
	if err != nil {
		config.Logger.Errorf("git commit error:%s:%s:%s", err, stderr.String(), out.String())
		exec.Command("git", "reset", "--hard").Run()
		return "", err
	}
	cmd = exec.Command("git", "push")
	cmd.Stdout, cmd.Stderr = &out, &stderr
	err = cmd.Run()
	if err != nil {
		config.Logger.Errorf("git push error:%s:%s", err, stderr.String())
		exec.Command("git", "reset", "--hard").Run()
		return "", err
	}
	return "success", nil
}
