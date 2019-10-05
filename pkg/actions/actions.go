package actions

import (
	"os"
	"strings"
)

//GetWorkspace returns workspace
func GetWorkspace() string {
	return os.Getenv("GITHUB_WORKSPACE")
}

//GetActor returns the actor
func GetActor() string {
	return os.Getenv("GITHUB_ACTOR")
}

//GetRepoOwner returns the repo owner where action is running
func GetRepoOwner() string {
	return strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[0]
}
