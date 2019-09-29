package store

import (
	"context"
	"net/http"

	"github.com/influxdata/influxdb-client-go"
	"github.com/rajatjindal/gfi-analytics-client/pkg/analytics"
)

//Influx stores data into influx db
type Influx struct {
	client *influxdb.Client
}

//New creates new influx client
func New(address, token string) (*Influx, error) {
	influx, err := influxdb.New(address, token, influxdb.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, err
	}

	return &Influx{
		client: influx,
	}, nil
}

func (i *Influx) Write(bucket, org string, data []*analytics.Entry) error {
	metrics := []influxdb.Metric{}

	for _, entry := range data {
		metrics = append(metrics, influxdb.NewRowMetric(
			map[string]interface{}{
				"click": 1,
			},
			"gfi-analytics",
			map[string]string{
				"owner":      entry.IssueData.Owner,
				"repo":       entry.IssueData.Repo,
				"issue-id":   entry.IssueData.IssueID,
				"user-agent": entry.BrowserInfo.UserAgent,
			},
			entry.BrowserInfo.ClickedAt,
		))
	}

	_, err := i.client.Write(context.TODO(), bucket, org, metrics...)
	return err
}
