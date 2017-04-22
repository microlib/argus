package main

import (
	"fmt"
	"github.com/microlib/argus/simple"
)

var (
	logger simple.Logger
	config simple.Config
)

func main() {

	logger.Level = "debug"
	logger.Info("Testing this all out")
	var cfg = config.Init("config.json")
	server, err := cfg.GetServer("haproxy-econdary")
	if err != nil {
		logger.Error(err.Error())
	}
	server, err = cfg.GetServer("haproxy-secondary")
	logger.Debug(fmt.Sprintf("level %s ", cfg.Level))
	logger.Level = cfg.Level
	logger.Trace(fmt.Sprintf("server %v ", server))
	logger.Info(fmt.Sprintf("cache %v ", cfg.Cache))

}
