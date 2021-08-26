package api

import (
	"fmt"
	"fssync/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/gin-gonic/gin"
)

// 手动同步接口
func GetFileFunc(c *gin.Context) {
	appName := c.Query("app")
	zone := c.Query("env")
	if appName == "" || zone == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": "false",
			"msg":    "Query value Missing",
		})
		return
	}
	dirpath := fmt.Sprintf("%s/env/%s/%s/", config.Config.RepoDir, zone, appName)
	filelist := []string{}
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			path = strings.ReplaceAll(path, dirpath, "")
			filelist = append(filelist, path)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "false",
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "true",
		"msg":    filelist,
	})
}
