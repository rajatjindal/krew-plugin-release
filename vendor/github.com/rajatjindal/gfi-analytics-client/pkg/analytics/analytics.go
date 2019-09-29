package analytics

import "time"

//Entry is one analytic entry
type Entry struct {
	IssueData   IssueData   `json:"issue-data"`
	BrowserInfo BrowserInfo `json:"browser-info"`
}

//BrowserInfo is details of browser
type BrowserInfo struct {
	UserAgent  string    `json:"user-agent"`
	RemoteAddr string    `json:"remote-addr"`
	ClickedAt  time.Time `json:"clicked-at"`
	RawURL     string    `json:"raw-url"`
}

//IssueData signifies IssueData
type IssueData struct {
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	IssueID string `json:"issue-id"`
}
