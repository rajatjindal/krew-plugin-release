package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/rajatjindal/krew-plugin-release/pkg/git"
	"github.com/rajatjindal/krew-plugin-release/pkg/krew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	ugit "gopkg.in/src-d/go-git.v4"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "tool to make PR to krew-plugin-release",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("reading release payload")
		releaseInfo, err := actions.GetReleaseInfo()
		if err != nil {
			logrus.Fatal(err)
		}

		dir, err := ioutil.TempDir("", "krew-index-")
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("will operate in tempdir %s", dir)
		repo, err := updateOriginFromUpstream(dir)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("creating branch %s", releaseInfo.GetTagName())
		err = git.CreateBranch(repo, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("checking out branch %s", releaseInfo.GetTagName())
		err = git.CheckoutBranch(repo, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("update plugin manifest with latest release info")
		err = krew.UpdatePluginManifest(dir, "modify-secret", releaseInfo)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("pushing changes to branch %s", releaseInfo.GetTagName())
		err = git.AddCommitAndPush(repo, "new version of modify-secret", releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("submitting the pr")
		err = submitPR(releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func stringp(s string) *string {
	return &s
}

func submitPR(branchName string) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("KREW_PLUGIN_RELEASE_TOKEN")})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)

	pr := &github.NewPullRequest{
		Title: stringp("release new version of modify-secret"),
		Head:  stringp(fmt.Sprintf("rajatjindal:%s", branchName)),
		Base:  stringp("master"),
		Body:  stringp("hey krew-index team, I would like to open this PR to release new version of modify-secret"),
	}

	_, _, err := client.PullRequests.Create(context.TODO(), "rajatjin", "krew-index", pr)
	if err != nil {
		return err
	}

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

func updateOriginFromUpstream(dir string) (*ugit.Repository, error) {
	logrus.Infof("Cloning %s", "https://github.com/rajatjindal/krew-index.git")
	repo, err := git.Clone("https://github.com/rajatjindal/krew-index.git", git.GetMasterBranchRefs(), dir)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Adding remote %s", "https://github.com/rajatjin/krew-index.git")
	remote, err := git.AddUpstream(repo, "https://github.com/rajatjin/krew-index.git")
	if err != nil {
		return nil, err
	}

	logrus.Info("fetching upstream")
	err = git.FetchUpstream(remote)
	if err != nil {
		return nil, err
	}

	logrus.Infof("pushing to origin/master of %s", "https://github.com/rajatjindal/krew-index.git")
	err = git.PushOriginMaster(repo)
	if err != nil && err.Error() != "already up-to-date" {
		return nil, err
	}

	return repo, nil
}
