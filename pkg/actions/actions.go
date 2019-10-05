package actions

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Inputs is action inputs
type Inputs struct {
	PluginName                string
	Token                     string
	TokenUserHandle           string
	TokenUserEmail            string
	TokenUserName             string
	UpstreamKrewIndexRepoName string
	UpstreamKrewIndexOwner    string
}

//Derived is derived data
type Derived struct {
	UpstreamCloneURL string
	LocalCloneURL    string
}

//ActionData is action data
type ActionData struct {
	Workspace   string
	Actor       string
	Repo        string
	RepoOwner   string
	Inputs      Inputs
	ReleaseInfo *github.RepositoryRelease
	Derived     Derived
}

//Action defines interface to get data from actions
type Action interface {
	GetActionData() ActionData
}

//RealAction is the real action
type RealAction struct{}

//GetActionData returns action data
func (r RealAction) GetActionData() (ActionData, error) {
	payload, err := r.getPayload()
	if err != nil {
		return ActionData{}, err
	}

	releaseInfo, err := getReleaseInfo(payload)
	if err != nil {
		return ActionData{}, err
	}

	tokenUserHandle := os.Getenv("KREW_PLUGIN_RELEASE_USER")
	if tokenUserHandle == "" {
		tokenUserHandle = strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[0]
	}

	token := os.Getenv("KREW_PLUGIN_RELEASE_TOKEN")
	tokenUserEmail, tokenUserName, err := r.getUserInfo(tokenUserHandle, token)
	if err != nil {
		return ActionData{}, err
	}

	upstreamKrewIndexRepoName := os.Getenv("UpstreamKrewIndexRepoName")
	if upstreamKrewIndexRepoName == "" {
		upstreamKrewIndexRepoName = "krew-index"
	}

	upstreamKrewIndexRepoOwner := os.Getenv("UpstreamKrewIndexRepoOwner")
	if upstreamKrewIndexRepoOwner == "" {
		upstreamKrewIndexRepoName = "kubernetes-sigs"
	}

	return ActionData{
		Workspace:   os.Getenv("GITHUB_WORKSPACE"),
		Actor:       os.Getenv("GITHUB_ACTOR"),
		Repo:        strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1],
		RepoOwner:   strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[0],
		ReleaseInfo: releaseInfo,
		Inputs: Inputs{
			PluginName:                r.getInputForAction("plugin-name"),
			TokenUserHandle:           tokenUserHandle,
			Token:                     os.Getenv("KREW_PLUGIN_RELEASE_TOKEN"),
			TokenUserEmail:            tokenUserEmail,
			TokenUserName:             tokenUserName,
			UpstreamKrewIndexOwner:    upstreamKrewIndexRepoOwner,
			UpstreamKrewIndexRepoName: upstreamKrewIndexRepoName,
		},
		Derived: Derived{
			UpstreamCloneURL: getRepoURL(upstreamKrewIndexRepoOwner, upstreamKrewIndexRepoName),
			LocalCloneURL:    getRepoURL(strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[10], strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1]),
		},
	}, nil
}

func getRepoURL(owner, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

//getInputForAction gets input to action
func (r RealAction) getInputForAction(key string) string {
	return os.Getenv(fmt.Sprintf("INPUT_%s", strings.ToUpper(key)))
}

//GetPayload reads payload and returns it
func (r RealAction) getPayload() ([]byte, error) {
	eventJSONPath := os.Getenv("GITHUB_EVENT_PATH")
	data, err := ioutil.ReadFile(eventJSONPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r RealAction) getUserInfo(username, token string) (string, string, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(context.TODO(), username)
	if err != nil {
		return "", "", err
	}

	return user.GetName(), user.GetEmail(), nil
}
