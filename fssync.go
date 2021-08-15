package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Getip(env, appName *string) ([]string, error) {
	return []string{"10.0.0.11"}, nil
}

func contains(s *[]string, str *string) bool {
	for _, v := range *s {
		if v == *str {
			return true
		}
	}
	return false
}

func SyncCron(srcPath, ip *string) {
	err := exec.Command("/usr/bin/rsync", "-avz", "--timeout=3", "--owner=wls81", "--group=wls", *srcPath, "root@"+*ip+":/").Run()
	if err != nil {
		Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", *srcPath, *ip, err.Error())
	} else {
		Logger.Infof("Rsync success:rsync -az %s* %s:/", *srcPath, *ip)
	}
}

func SyncHttp(srcPath, ip *string, ch chan string) {
	out, err := exec.Command("/usr/bin/rsync", "-avz", "--timeout=3", "--owner=wls81", "--group=wls", *srcPath, "root@"+*ip+":/").Output()
	if err != nil {
		Logger.Infof("Rsync error:%s", out)
		Logger.Errorf("Rsync error:[%s]-[%s]-[%s]", *srcPath, *ip, err.Error())
		ch <- *ip + ":" + string(out)
	} else {
		Logger.Infof("Rsync success:rsync -az %s* %s:/", *srcPath, *ip)
		ch <- *ip + ":success"
	}
}

// 获取与上个版本有差异的文件列表
func GetDiffFile() ([]string, error) {
	os.Chdir(repoDir)
	cmd := exec.Command("git", "pull", "origin", "main")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		Logger.Errorf("git pull error:%s:%s", err, stderr.String())
		return nil, err
	} else {
		Logger.Info("git pull: " + strings.Trim(out.String(), "\n"))
		if !strings.Contains(out.String(), "Already up to date") {
			// git log -p -1 --oneline
			out, _ := exec.Command("git", "diff", "head^", "--name-only").Output()
			files := strings.Split(strings.Trim(string(out), "\n"), "\n")
			return files, nil
		} else {
			Logger.Infof("No diff file")
			return nil, nil
		}
	}
}

func Worker() {
	var src []string
	diffFiles, err := GetDiffFile()
	if err != nil {
		Logger.Error("Get diff file error:", err.Error())
	} else if diffFiles != nil {
		Logger.Info("get files:", diffFiles)
	}
	for _, file := range diffFiles {
		Logger.Infof("Start sync file:%s", file)
		dirs := strings.Split(file, "/")
		if len(dirs) > 3 {
			mod, env, appName := dirs[0], dirs[1], dirs[2]
			srcPath := fmt.Sprintf("%s/%s/%s/%s/", repoDir, mod, env, appName)

			if !contains(&src, &srcPath) {
				src = append(src, srcPath)
				ipList, err := Getip(&env, &appName)
				if err != nil {
					Logger.Error("Get ip error:", err.Error())
					return
				}
				Logger.Infof("[Sync info]src:%s;mod:%s;env:%s;app:%s;iplist:%s", srcPath, mod, env, appName, ipList)
				for _, ip := range ipList {
					go SyncCron(&srcPath, &ip)
				}
			}
		}
	}
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: logfile,
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

const (
	workdir  = "/Users/pangru/Documents/golang/fssync/"
	logfile  = "/Users/pangru/Documents/golang/fssync/syncer.log"
	repoDir  = "/Users/pangru/Documents/golang/fssync/serverfiles" // git本地仓库目录
	interval = "10"                                                //cron 间隔
	timeout  = "3"
	owner    = "wls81"
	group    = "wls"
)

var (
	Logger *zap.SugaredLogger
)

func Cron() {
	cron := cron.New()

	cron.AddFunc(fmt.Sprintf("@every %ss", interval), func() {
		Logger.Infof("start cron")
		Worker()
		Logger.Infof("end cron")
	})

	cron.Start()
}

func syncApp(c *gin.Context) {
	env := c.DefaultQuery("env", "all")
	appName := c.Query("app")
	mod := c.Query("mod")
	ipList, err := Getip(&env, &appName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "get ip from cmdb failed"})
	}
	srcPath := fmt.Sprintf("%s/%s/%s/%s/", repoDir, mod, env, appName)
	Logger.Infof("[Sync info]mod:%s;env:%s;app:%s;iplist:%s", mod, env, appName, ipList)

	ch := make(chan string, len(ipList))
	for _, ip := range ipList {
		go SyncHttp(&srcPath, &ip, ch)
	}
	var resp []string
	for i := 0; i < len(ipList); i++ {
		r := <-ch
		resp = append(resp, r)
	}
	c.IndentedJSON(http.StatusOK, resp)
}
func main() {
	InitLogger()
	// Worker()
	Cron()
	router := gin.Default()
	router.GET("/sync", syncApp)

	router.Run(":8080")

}
