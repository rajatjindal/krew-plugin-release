package actions

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//ActionData defines interface to get data from actions
type ActionData interface {
	GetWorkspace() string
	GetActor() string
	GetRepoOwner() string
	GetInputForAction(key string) string
	GetPluginName() string
	GetPayload() ([]byte, error)
}

//RealAction is the real action
type RealAction struct{}

//GetWorkspace returns workspace
func (r RealAction) GetWorkspace() string {
	return os.Getenv("GITHUB_WORKSPACE")
}

//GetActor returns the actor
func (r RealAction) GetActor() string {
	return os.Getenv("GITHUB_ACTOR")
}

//GetRepoOwner returns the repo owner where action is running
func (r RealAction) GetRepoOwner() string {
	return strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[0]
}

//GetPluginName gets the plugin name
func (r RealAction) GetPluginName() string {
	return r.GetInputForAction("plugin-name")
}

//GetInputForAction gets input to action
func (r RealAction) GetInputForAction(key string) string {
	return os.Getenv(fmt.Sprintf("INPUT_%s", strings.ToUpper(key)))
}

//GetPayload reads payload and returns it
func (r RealAction) GetPayload() ([]byte, error) {
	eventJSONPath := os.Getenv("GITHUB_EVENT_PATH")
	data, err := ioutil.ReadFile(eventJSONPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
