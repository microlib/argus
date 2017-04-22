package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Response struct {
	Status  string      `json:"status"`
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

var (
	deploy      string = "deploy"
	undeploy    string = "deploy"
	appexecute  string = "appexecute"
	serverstats string = "serverstats"
)

func DeployRepo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
	errJson := json.Unmarshal(body, &f)
	if errJson != nil {
		logger.Error(err.Error())
	}

	m := f.(map[string]interface{})
	var response Response

	validateErr := validateParameters(deploy, m)
	if validateErr != nil {
		logger.Error(err.Error())
		response = Response{Status: "401", Result: "ko", Message: "ERROR Bad request " + err.Error(), Payload: nil}

	} else {

		logger.Debug("Directory " + config.Basedir + "/" + m["Application"].(string))
		logger.Debug("Repo " + m["Repo"].(string))

		// check the cache first
		//val, err := client.Get(m["application"].(string)).Result()
		//if err != nil {
		//	logger.Error(err.Error())
		//	response = Response{Status: "500", Result: "ko", Message: "ERROR Application " + m["application"].(string) + " redis cache ", Payload: nil}
		//} else {
		//	if val != "" {
		//		response = Response{Status: "200", Result: "ok", Message: "Application " + m["application"].(string) + " already deployed use force true to override", Payload: nil}
		//	} else {
		giterr := GetRepo(config.Basedir+"/"+m["Application"].(string), m["BaseRepoUrl"].(string)+"/"+m["Repo"].(string))
		if giterr != nil {
			logger.Error(giterr.Error())
			response = Response{Status: "500", Result: "ko", Message: "ERROR Application " + m["Application"].(string) + " failed to deploy", Payload: nil}
		} else {
			response = Response{Status: "200", Result: "ok", Message: "Application " + m["Application"].(string) + " succesfully deployed", Payload: nil}
			file, _ := GetAppConfig(config.Basedir + "/" + m["Application"].(string) + "/config.json")
			logger.Debug(string(file))
			//client.Set(m["application"].(string), file, 0)
		}
	}
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func UndeployRepo(w http.ResponseWriter, r *http.Request) {
	var (
		f        interface{}
		response Response
	)

	body, err := ioutil.ReadAll(r.Body)
	errJson := json.Unmarshal(body, &f)
	if errJson != nil {
		logger.Error(err.Error())
	}

	m := f.(map[string]interface{})
	validateErr := validateParameters(undeploy, m)
	if validateErr != nil {
		logger.Error(err.Error())
		response = Response{Status: "401", Result: "ko", Message: "ERROR Bad request " + err.Error(), Payload: nil}

	} else {
		giterr := RemoveContents(config.Basedir + "/" + m["Application"].(string))
		if giterr != nil {
			logger.Error(giterr.Error())
			response = Response{Status: "500", Result: "ko", Message: "ERROR Application " + m["Application"].(string) + " not removed", Payload: nil}
		} else {
			response = Response{Status: "200", Result: "ok", Message: "Application " + m["Application"].(string) + " succesfully removed", Payload: nil}
			//client.Del(m["Application"].(string))
		}
	}

	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func GetServerStats(w http.ResponseWriter, r *http.Request) {
	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
	errJson := json.Unmarshal(body, &f)
	if errJson != nil {
		logger.Error(err.Error())
	}
	var response Response
	idle, total := GetCPUSample()
	response = Response{Status: "200", Result: "ok", Message: "CPU total=" + strconv.FormatUint(total, 10) + " idle=" + strconv.FormatUint(idle, 10), Payload: nil}
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func CallExecute(w http.ResponseWriter, r *http.Request) {
	var (
		f        interface{}
		response Response
	)

	body, err := ioutil.ReadAll(r.Body)
	errJson := json.Unmarshal(body, &f)
	if errJson != nil {
		logger.Error(err.Error())
	}

	m := f.(map[string]interface{})
	validateErr := validateParameters(undeploy, m)
	if validateErr != nil {
		logger.Error(err.Error())
		response = Response{Status: "401", Result: "ko", Message: "ERROR Bad request " + err.Error(), Payload: nil}

	} else {
		output := AppExecute(config.Basedir+"/"+m["Application"].(string), m["Action"].(string))
		response = Response{Status: "200", Result: "ok", Message: "Application " + m["Application"].(string) + " response " + output, Payload: nil}
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, string(b))
	}
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func validateParameters(handlerType string, mp map[string]interface{}) error {

	switch handlerType {
	case deploy:
		if mp["Application"] == nil || mp["Application"].(string) == "" {
			return errors.New("missing 'Application' parameter")
		} else if mp["BaseRepoUrl"] == nil || mp["BaseRepoUrl"].(string) == "" {
			return errors.New("missing 'BaseRepoUrl' parameter")
		} else if mp["Repo"] == nil || mp["Repo"].(string) == "" {
			return errors.New("missing 'Repo' parameter")
		} else if mp["Action"] == nil || mp["Action"].(string) == "" {
			return errors.New("missing 'Action' parameter")
		} else if mp["Token"] == nil || mp["Token"].(string) == "" {
			return errors.New("missing 'Token' parameter")
		} else {
			return nil
		}
	case undeploy:
		if mp["Application"] == nil || mp["Application"].(string) == "" {
			return errors.New("missing 'Application' parameter")
		} else if mp["Token"] == nil || mp["Token"].(string) == "" {
			return errors.New("missing 'Token' parameter")
		} else {
			return nil
		}
	case serverstats, appexecute:
		if mp["Application"] == nil || mp["Application"].(string) == "" {
			return errors.New("missing 'Application' parameter")
		} else if mp["Action"] == nil || mp["Action"].(string) == "" {
			return errors.New("missing 'Action' parameter")
		} else if mp["Token"] == nil || mp["Token"].(string) == "" {
			return errors.New("missing 'Token' parameter")
		} else {
			return nil
		}
	default:
		return nil
	}

}
