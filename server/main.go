package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"../simple"
	"github.com/go-redis/redis"
)

var (
	logger simple.Logger
	client *redis.Client
	config simple.Config
)

func startHttpServer(port string) *http.Server {
	srv := &http.Server{Addr: ":" + port}

	// add all the routes and link to handlers
	http.HandleFunc("/deploy", DeployRepo)
	http.HandleFunc("/undeploy", UndeployRepo)
	http.HandleFunc("/serverstats", GetServerStats)
	http.HandleFunc("/appexecute", CallExecute)
	http.HandleFunc("/isalive", IsAlive)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("Httpserver: ListenAndServe() error: %s", err.Error()))
		}
	}()

	return srv
}

func cleanup() {
	logger.Info("main: server cleanup")
}

func main() {
	var cfg simple.Config
	config = cfg.Init("config.json")

	logger.Info(fmt.Sprintf("Debug LMZ %v", config.Cache))
	logger.Level = config.Level

	if config.Cache == "true" {
		server, err := config.GetServer("redis")
		if err != nil {
			logger.Error(err.Error())
		}

		client = redis.NewClient(&redis.Options{
			Addr:     server.Host + ":" + server.Port,
			Password: "",
			DB:       0,
		})

		logger.Info(client.Ping().String())
		logger.Info(fmt.Sprintf("Redis info %s %s", server.Host, server.Port))
	}

	srv := startHttpServer(config.Port)
	logger.Info(fmt.Sprintf("main: starting server on port %s", srv.Addr))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	cleanup()
	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("main: server shutdown successfully")
	os.Exit(code)
}
