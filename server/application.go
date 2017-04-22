package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func AppExecute(app, action string) string {
	var output string
	switch action {
	case "start":
		output = start(app)
	case "stop":
		output = stop(app)
	case "status":
		output = status(app)
	case "build":
		output = build(app)

	default:
		output = "unrecognized action"
	}
	return output
}

// This can be drastically refactored as the only change is the action
// I kept this "as is" incase we need to some other elaboration
func start(app string) string {
	if ok := isDeployed(app); ok {
		logger.Trace("Application " + app)
		out, _ := exec.Command("sh", "-c", app+"/script.sh start "+app).Output()
		logger.Debug("Response from script [start] " + string(out))
		return string(out)
	} else {
		return app + " not deployed"
	}
}

func stop(app string) string {
	if ok := isDeployed(app); ok {
		out, _ := exec.Command("sh", "-c", app+"/script.sh stop").Output()
		return string(out)
	} else {
		return app + " not deployed"
	}
}

func status(app string) string {
	if ok := isDeployed(app); ok {
		out, _ := exec.Command("sh", "-c", app+"/script.sh status").Output()
		return string(out)
	} else {
		return app + " not deployed"
	}
}

func build(app string) string {
	if ok := isDeployed(app); ok {
		logger.Trace("Application " + app)
		out, _ := exec.Command("sh", "-c", app+"/script.sh build "+app).Output()
		logger.Debug("Response from script [build] " + string(out))
		return string(out)
	} else {
		return app + " not deployed"
	}
}

func isDeployed(app string) bool {
	if _, err := os.Stat(app); err != nil {
		if os.IsNotExist(err) {
		} else {
			// other error
		}
		return false
	} else {
		return true
	}
}

// execCmd with wait group
func exeCmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}
