package api

import (
	"fmt"
	"fssync/config"
	"fssync/core"
	"fssync/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncFunc(c *gin.Context) {
	env := c.DefaultQuery("env", "all")
	appName := c.Query("app")
	mod := c.Query("mod")
	ipList, err := util.Getip(&env, &appName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "get ip from cmdb failed"})
	}
	srcPath := fmt.Sprintf("%s/%s/%s/%s/", config.RepoDir, mod, env, appName)
	config.Logger.Infof("[Sync info]mod:%s;env:%s;app:%s;iplist:%s", mod, env, appName, ipList)

	ch := make(chan string, len(ipList))
	for _, ip := range ipList {
		go core.SyncHttp(&srcPath, &ip, ch)
	}
	var resp []string
	for i := 0; i < len(ipList); i++ {
		r := <-ch
		resp = append(resp, r)
	}
	c.IndentedJSON(http.StatusOK, resp)
}
