package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/rajatjindal/krew-plugin-release/pkg/git"
	"github.com/rajatjindal/krew-plugin-release/pkg/krew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	ugit "gopkg.in/src-d/go-git.v4"
)

const prBody = `hey krew-index team,

I would like to open this PR to release new version %s of %s on behalf of %s.

Thanks,
[krew-plugin-release](https://github.com/rajatjindal/krew-plugin-release)`

type actionData struct {
	pluginName             string
	tagName                string
	token                  string
	localKrewIndexOwner    string
	localKrewIndexRepo     string
	localRemoteName        string
	upstreamKrewIndexRepo  string
	upstreamKrewIndexOwner string
	upstreamRemoteName     string
	actor                  string
	branchName             string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "github action to open PR for krew-index on release of new version of krew-plugin",
	Run: func(cmd *cobra.Command, args []string) {
		action := actions.RealAction{}

		data := actionData{
			pluginName:             action.GetPluginName(),
			token:                  os.Getenv("KREW_PLUGIN_RELEASE_TOKEN"),
			upstreamKrewIndexOwner: action.GetInputForAction("upstream-krew-index-owner"),
			localKrewIndexOwner:    action.GetRepoOwner(),
			localKrewIndexRepo:     fmt.Sprintf("https://github.com/%s/krew-index.git", action.GetRepoOwner()),
			upstreamKrewIndexRepo:  fmt.Sprintf("https://github.com/%s/krew-index.git", action.GetInputForAction("upstream-krew-index-owner")),
			localRemoteName:        "local",
			upstreamRemoteName:     "upstream",
			actor:                  action.GetActor(),
		}

		logrus.Info("reading release payload")
		releaseInfo, err := actions.GetReleaseInfo(action)
		if err != nil {
			logrus.Fatal(err)
		}

		if releaseInfo.GetPrerelease() {
			logrus.Infof("%s is a pre-release. not opening the PR", releaseInfo.GetTagName())
			logrus.Exit(0)
		}
		data.tagName = releaseInfo.GetTagName()
		data.branchName = releaseInfo.GetTagName()

		dir, err := ioutil.TempDir("", "krew-index-")
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("will operate in tempdir %s", dir)
		repo, err := cloneRepos(data, dir)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("creating branch %s", releaseInfo.GetTagName())
		err = git.CreateBranch(repo, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("update plugin manifest with latest release info")

		templateFile := filepath.Join(action.GetWorkspace(), ".krew.yaml")
		actualFile := filepath.Join(dir, "plugins", krew.PluginFileName(data.pluginName))
		err = krew.UpdatePluginManifest(templateFile, actualFile, releaseInfo)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("pushing changes to branch %s", data.tagName)
		err = git.AddCommitAndPush(repo, fmt.Sprintf("new version of %s", data.pluginName), data.localRemoteName, data.tagName)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("submitting the pr")
		err = submitPR(data)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func stringp(s string) *string {
	return &s
}

func submitPR(data actionData) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("KREW_PLUGIN_RELEASE_TOKEN")})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)

	prr := &github.NewPullRequest{
		Title: stringp(fmt.Sprintf("release new version %s of %s", data.tagName, data.pluginName)),
		Head:  stringp(fmt.Sprintf("%s:%s", data.localKrewIndexOwner, data.branchName)),
		Base:  stringp("master"),
		Body:  stringp(fmt.Sprintf(prBody, data.tagName, data.pluginName, data.actor)),
	}

	pr, _, err := client.PullRequests.Create(context.TODO(), data.upstreamKrewIndexOwner, "krew-index", prr)
	if err != nil {
		return err
	}

	logrus.Infof("pr %q opened for releasing new version", pr.GetHTMLURL())
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cloneRepos(data actionData, dir string) (*ugit.Repository, error) {
	logrus.Infof("Cloning %s", data.upstreamKrewIndexRepo)
	repo, err := git.Clone(data.upstreamKrewIndexRepo, data.upstreamRemoteName, git.GetMasterBranchRefs(), dir)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Adding remote %s at %s", data.localRemoteName, data.localKrewIndexRepo)
	_, err = git.AddUpstream(repo, data.localRemoteName, data.localKrewIndexRepo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
