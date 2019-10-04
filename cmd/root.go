package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/rajatjindal/krew-plugin-release/pkg/git"
	"github.com/rajatjindal/krew-plugin-release/pkg/krew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "tool to make PR to krew-plugin-release",
	Run: func(cmd *cobra.Command, args []string) {
		releaseInfo, err := actions.GetReleaseInfo()
		if err != nil {
			logrus.Fatal(err)
		}

		dir, err := ioutil.TempDir("", "krew-index-")
		if err != nil {
			logrus.Fatal(err)
		}

		err = UpdateOriginFromUpstream(dir)
		if err != nil {
			logrus.Fatal(err)
		}

		err = git.CreateBranch(dir, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		err = git.CheckoutBranch(dir, releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}

		err = krew.UpdatePluginManifest(dir, "modify-secret", releaseInfo)
		if err != nil {
			logrus.Fatal(err)
		}

		err = git.CommitAndPush(dir, "new version of modify-secret", releaseInfo.GetTagName())
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func stringp(s string) *string {
	return &s
}

func submitPR() {
	// pr := github.NewPullRequest{
	// 	Title: stringp("release new version of modify-secret"),
	// 	Head:  stringp("rajatjindal:krew-index"),
	// 	Base:  stringp("master"),
	// 	Body:  stringp("hey krew-index team, I would like to open this PR to release new version of modify-secret"),
	// }
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func UpdateOriginFromUpstream(dir string) error {
	err := git.Clone("https://github.com/rajatjindal/krew-index.git", git.GetMasterBranchRefs(), dir)
	if err != nil {
		return err
	}

	err = git.AddUpstream(dir, "https://github.com/rajatjin/krew-index.git")
	if err != nil {
		return err
	}

	err = git.FetchUpstream(dir)
	if err != nil {
		return err
	}

	err = git.RebaseUpstream(dir)
	if err != nil {
		return err
	}

	err = git.PushOriginMaster(dir)
	if err != nil {
		return err
	}

	return nil
}
