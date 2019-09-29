package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/rajatjindal/gfi-analytics-client/pkg/analytics"
	"github.com/rajatjindal/gfi-analytics-client/pkg/store"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gfi-analytics-client",
	Short: "client to fetch data from gfi-analytics",

	Run: func(cmd *cobra.Command, args []string) {
		idb, err := store.New("http://mypi.local:8086", "")
		if err != nil {
			logrus.Fatal("1 ", err)
		}

		data, err := pullAnalytics()
		if err != nil {
			logrus.Fatal("2 ", err)
		}

		err = idb.Write("gfi-analytics", "goodfirstissue", data)
		if err != nil {
			logrus.Fatal("3 ", err)
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

func pullAnalytics() ([]*analytics.Entry, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "rajatjindal.o6s.io",
		Path:     "/gfi-analytics",
		RawQuery: "get-analytics=true",
	}

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := []*analytics.Entry{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
