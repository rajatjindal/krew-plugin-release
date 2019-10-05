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

type actionInputs struct {
	PluginName             string
	Token                  string
	UpstreamKrewIndexOwner string
	LocalKrewIndexOwner    string
	localKrewIndexRepo     string
	upstreamKrewIndexRepo  string
	localRemoteName        string
	upstreamRemoteName     string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "tool to make PR to krew-plugin-release",
	Run: func(cmd *cobra.Command, args []string) {
		action := actions.RealAction{}

		inputs := actionInputs{
			PluginName:             action.GetPluginName(),
			Token:                  os.Getenv("KREW_PLUGIN_RELEASE_TOKEN"),
			UpstreamKrewIndexOwner: action.GetInputForAction("upstream-krew-index-owner"),
			LocalKrewIndexOwner:    action.GetRepoOwner(),
			localKrewIndexRepo:     fmt.Sprintf("https://github.com/%s/krew-index.git", action.GetRepoOwner()),
			upstreamKrewIndexRepo:  fmt.Sprintf("https://github.com/%s/krew-index.git", action.GetInputForAction("upstream-krew-index-owner")),
			localRemoteName:        "local",
			upstreamRemoteName:     "upstream",
		}

		logrus.Info("reading release payload")
		releaseInfo, err := actions.GetReleaseInfo(action)
		if err != nil {
			logrus.Fatal(err)
		}

		dir, err := ioutil.TempDir("", "krew-index-")
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("will operate in tempdir %s", dir)
		repo, err := cloneRepos(inputs, dir)
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
		actualFile := filepath.Join(dir, "plugins", krew.PluginFileName(inputs.PluginName))
		err = krew.UpdatePluginManifest(templateFile, actualFile, releaseInfo)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("pushing changes to branch %s", releaseInfo.GetTagName())
		err = git.AddCommitAndPush(repo, fmt.Sprintf("new version of %s", inputs.PluginName), inputs.localRemoteName, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("submitting the pr")
		err = submitPR(inputs, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func stringp(s string) *string {
	return &s
}

func submitPR(inputs actionInputs, branchName string) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("KREW_PLUGIN_RELEASE_TOKEN")})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)

	prr := &github.NewPullRequest{
		Title: stringp(fmt.Sprintf("release new version of %s", inputs.PluginName)),
		Head:  stringp(fmt.Sprintf("%s:%s", inputs.LocalKrewIndexOwner, branchName)),
		Base:  stringp("master"),
		Body:  stringp("hey krew-index team, I would like to open this PR to release new version of modify-secret"),
	}

	pr, _, err := client.PullRequests.Create(context.TODO(), inputs.UpstreamKrewIndexOwner, "krew-index", prr)
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

func cloneRepos(inputs actionInputs, dir string) (*ugit.Repository, error) {
	logrus.Infof("Cloning %s", inputs.upstreamKrewIndexRepo)
	repo, err := git.Clone(inputs.upstreamKrewIndexRepo, inputs.upstreamRemoteName, git.GetMasterBranchRefs(), dir)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Adding remote %s at %s", inputs.localRemoteName, inputs.localKrewIndexRepo)
	_, err = git.AddUpstream(repo, inputs.localRemoteName, inputs.localKrewIndexRepo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
