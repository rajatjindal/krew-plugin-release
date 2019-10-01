package cmd

import (
	"fmt"
	"os"

	"github.com/rajatjindal/krew-plugin-release-go/pkg/actions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krew-plugin-release",
	Short: "tool to make PR to krew-plugin-release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("from inside golang")
		_, err := actions.GetReleaseInfo()
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
