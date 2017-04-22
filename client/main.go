package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/microlib/argus/simple"
)

type Payload struct {
	Application string
	Token       string
	Server      string
	BaseRepoUrl string
	Repo        string
	Refs        string
	Action      string
	Other       interface{}
}

func main() {

	// Endpoints
	// to deploy an app /deploy
	// to remove an app /undeploy
	// get server status /serverstats
	// execute an app script /appexecute actions: build,start,stop,status

	var (
		cfg, config simple.Config
		logger      simple.Logger
		client      *redis.Client
		server      simple.Server
	)

	config = cfg.Init("../server/config.json")
	logger.Level = config.Level

	// If we have cache set to true
	if config.Cache == "true" {
		var err error
		server, err = config.GetServer("redis")
		if err != nil {
			logger.Error(err.Error())
		}

		client = redis.NewClient(&redis.Options{
			Addr:     server.Host + ":" + server.Port,
			Password: "",
			DB:       0,
		})

		value, err := client.Ping().Result()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(-1)
		}
		logger.Info("client " + client.String())
		logger.Info("response from redis server " + value)
	}

	// Get the list of servers and deploy to the first server in the list
	// If this fails try the second ... to the last server in the list
	// If all fails recommend adding a new 'server'

	var payload Payload

	// cli format
	// ./client <token> <application> <url> <action> <base repo url>
	// ./client 3213432543 go-redis-service deploy github.com/dimitraz

	for _, server := range config.Servers {
		if strings.Contains(server.Name, "application") {
			url := "http://" + server.Host + ":" + server.Port + "/" + os.Args[3]
			logger.Info("URL:>" + url)
			logger.Info(fmt.Sprintf("Args %v ", os.Args))

			payload = Payload{Token: os.Args[1], Application: os.Args[2], Action: os.Args[4], Server: url, BaseRepoUrl: os.Args[5], Repo: os.Args[2], Refs: "master", Other: nil}
			b, _ := json.Marshal(payload)
			logger.Info(string(b))

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
			req.Header.Set("X-Custom-Header", "myvalue")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			logger.Info("response Status:" + resp.Status)
			logger.Info("response Headers: ")
			for k, _ := range resp.Header {
				logger.Info("  key:" + string(k))
				//logger.Info("  value:" + string(v[k]))
			}
			body, _ := ioutil.ReadAll(resp.Body)
			logger.Info("response Body:" + string(body))
			// if this is ok then break
			// else move onto the next available server
			if resp.StatusCode == 200 {
				break
			}
		}
	}

	//if index == servers.length {
	//	logger.Info("No server to deploy to (cpu and mem resources are exhuasted) please consider addin a new server")
	//:w
}
