package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Para struct {
	Workdir  string `json:"workdir" yaml:"workdir"`
	Logfile  string `json:"logfile" yaml:"logfile"`
	RepoDir  string `json:"repodir" yaml:"repodir"`
	Easypub  string `json:"easypub" yaml:"easypub"`
	Interval string `json:"interval" yaml:"interval"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	Owner    string `json:"owner" yaml:"owner"`
	Group    string `json:"group" yaml:"group"`
	Env      string `json:"env" yaml:"env"`
}

var (
	Logger *zap.SugaredLogger
	VP     *viper.Viper
	Config Para
)
