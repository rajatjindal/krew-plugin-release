package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	ugit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func cloneRepos(actionData actions.ActionData, dir string) (*ugit.Repository, error) {
	logrus.Infof("Cloning %s", actionData.Derived.UpstreamCloneURL)
	repo, err := ugit.PlainClone(dir, false, &ugit.CloneOptions{
		URL:           actionData.Derived.UpstreamCloneURL,
		Progress:      os.Stdout,
		ReferenceName: plumbing.Master,
		SingleBranch:  true,
		Auth:          getAuth(actionData),
		RemoteName:    originNameUpstream,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("Adding remote %s at %s", originNameLocal, actionData.Derived.LocalCloneURL)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: originNameLocal,
		URLs: []string{actionData.Derived.LocalCloneURL},
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("creating branch %s", actionData.ReleaseInfo.GetTagName())
	err = createBranch(repo, actionData.ReleaseInfo.GetTagName())
	if err != nil {
		return nil, err
	}

	return repo, nil
}

//createBranch creates branch
func createBranch(repo *ugit.Repository, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// First try to create branch
	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: plumbing.NewBranchReferenceName(branchName),
	})

	if err == nil {
		return nil
	}

	//may be it already exists
	return w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
}

type commit struct {
	msg        string
	remoteName string
}

//addCommitAndPush commits and push
func addCommitAndPush(repo *ugit.Repository, commit commit, actionData actions.ActionData) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	w.Add(".")
	_, err = w.Commit(commit.msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  actionData.Inputs.TokenUserName,
			Email: actionData.Inputs.TokenUserEmail,
			When:  time.Now(),
		},
	})

	return repo.Push(&ugit.PushOptions{
		RemoteName: commit.remoteName,
		RefSpecs:   []config.RefSpec{config.DefaultPushRefSpec},
		Auth:       getAuth(actionData),
	})
}

func submitPR(actionData actions.ActionData) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: actionData.Inputs.Token})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)

	prr := &github.NewPullRequest{
		Title: getTitle(actionData),
		Head:  getHead(actionData),
		Base:  github.String("master"),
		Body:  getPRBody(actionData),
	}

	logrus.Infof("creating pr with title %q, \nhead %q, \nbase %q, \nbody %q",
		github.Stringify(getTitle(actionData)),
		github.Stringify(getHead(actionData)),
		"master",
		github.Stringify(getPRBody(actionData)),
	)

	pr, _, err := client.PullRequests.Create(
		context.TODO(),
		actionData.Inputs.UpstreamKrewIndexOwner,
		actionData.Inputs.UpstreamKrewIndexRepoName,
		prr,
	)
	if err != nil {
		return err
	}

	logrus.Infof("pr %q opened for releasing new version", pr.GetHTMLURL())
	return nil
}

func getTitle(actionData actions.ActionData) *string {
	s := fmt.Sprintf(
		"release new version %s of %s",
		actionData.ReleaseInfo.GetTagName(),
		actionData.Inputs.PluginName,
	)

	return github.String(s)
}

func getHead(actionData actions.ActionData) *string {
	s := fmt.Sprintf("%s:%s", actionData.RepoOwner, actionData.ReleaseInfo.GetTagName())
	return github.String(s)
}

func getPRBody(actionData actions.ActionData) *string {
	prBody := `hey krew-index team,

I would like to open this PR to publish version %s of %s on behalf of [%s](https://github.com/%s).

Thanks,
[krew-plugin-release](https://github.com/rajatjindal/krew-plugin-release)`

	s := fmt.Sprintf(prBody,
		fmt.Sprintf("`%s`", actionData.ReleaseInfo.GetTagName()),
		fmt.Sprintf("`%s`", actionData.Inputs.PluginName),
		actionData.Actor,
		actionData.Actor,
	)

	return github.String(s)
}

func getAuth(actionData actions.ActionData) transport.AuthMethod {
	return &githttp.BasicAuth{
		Username: actionData.Inputs.TokenUserHandle,
		Password: actionData.Inputs.Token,
	}
}
