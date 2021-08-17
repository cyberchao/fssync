package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Para struct {
	Workdir  string `json:"workdir" yaml:"workdir"`
	Logfile  string `json:"logfile" yaml:"logfile"`
	RepoDir  string `json:"repodir" yaml:"repodir"`
	Interval string `json:"interval" yaml:"interval"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	Owner    string `json:"owner" yaml:"owner"`
	Group    string `json:"group" yaml:"group"`
	Env      string `json:"env" yaml:"env"`
}

func Viper(path string) *viper.Viper {

	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&Config); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&Config); err != nil {
		fmt.Println(err)
	}
	return v
}
