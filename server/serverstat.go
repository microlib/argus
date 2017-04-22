package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					logger.Error(fmt.Sprintf("Error: %d %s %s", i, fields[i], err.Error()))
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}

func GetMEMSample() (total, free uint64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					logger.Error(fmt.Sprintf("Error: %d %s %s", i, fields[i], err.Error()))
				}
				total += val
				if i == 4 {
					free = val
				}
			}
			return
		}
	}
	return
}
