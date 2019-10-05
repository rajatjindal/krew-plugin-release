package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/rajatjindal/krew-plugin-release/pkg/krew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	originNameUpstream = "upstream"
	originNameLocal    = "local"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "github action to open PR for krew-index on release of new version of krew-plugin",
	Run: func(cmd *cobra.Command, args []string) {
		action := actions.LocalAction{}
		actionData, err := action.GetActionData()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info(actionData.Inputs, actionData.Derived)
		if actionData.ReleaseInfo.GetPrerelease() {
			logrus.Infof("%s is a pre-release. not opening the PR", actionData.ReleaseInfo.GetTagName())
			logrus.Exit(0)
		}

		tempdir, err := ioutil.TempDir("", "krew-index-")
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("will operate in tempdir %s", tempdir)
		repo, err := cloneRepos(actionData, tempdir)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("update plugin manifest with latest release info")

		templateFile := filepath.Join(actionData.Workspace, ".krew.yaml")
		actualFile := filepath.Join(tempdir, "plugins", krew.PluginFileName(actionData.Inputs.PluginName))
		err = krew.UpdatePluginManifest(templateFile, actualFile, actionData.ReleaseInfo)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("pushing changes to branch %s", actionData.ReleaseInfo.GetTagName())
		commit := commit{
			msg:        fmt.Sprintf("new version %s of %s", actionData.ReleaseInfo.GetTagName(), actionData.Inputs.PluginName),
			remoteName: originNameLocal,
		}
		err = addCommitAndPush(repo, commit, actionData)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("submitting the pr")
		err = submitPR(actionData)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
