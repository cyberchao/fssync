package api

import (
	"crypto/md5"
	"fmt"
	"fssync/config"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 手动同步接口
func GetFileFunc(c *gin.Context) {
	appName := c.Query("app")
	zone := c.Query("zone")
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

func Getk8sFileFunc(c *gin.Context) {
	appName := c.Query("app")
	zone := c.Query("env")
	if appName == "" || zone == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": "false",
			"msg":    "Query value Missing",
		})
		return
	}
	dirpath := fmt.Sprintf("%s/env/%s/%s", config.Config.RepoDir, zone, appName)
	filemap := make(map[string]string)
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			f, _ := ioutil.ReadFile(path)
			md5 := md5.Sum([]byte(f))

			key := strings.ReplaceAll(path, dirpath, "")
			value := fmt.Sprintf("%x_#_%s/%s%s", md5, zone, appName, key)
			filemap[key] = value
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
	commonpath := fmt.Sprintf("%s/env/%s/%s", config.Config.RepoDir, zone, "common_files")
	err = filepath.Walk(commonpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			f, _ := ioutil.ReadFile(path)
			md5 := md5.Sum([]byte(f))

			key := strings.ReplaceAll(path, dirpath, "")
			value := fmt.Sprintf("%x_#_%s/%s%s", md5, zone, appName, key)
			filemap[key] = value
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

	c.JSON(http.StatusOK, filemap)
}
