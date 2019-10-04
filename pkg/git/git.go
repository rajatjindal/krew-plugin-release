package git

import (
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"gopkg.in/src-d/go-git.v4"
	ugit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func getAuth() transport.AuthMethod {
	return &githttp.BasicAuth{
		Username: "rjindal",
		Password: os.Getenv("KREW_PLUGIN_RELEASE_TOKEN"),
	}
}

//Repository represents upstream repo
type Repository ugit.Repository

//Clone clones the repo
func Clone(origin, branch, dir string) (*ugit.Repository, error) {

	return ugit.PlainClone(dir, false, &ugit.CloneOptions{
		URL:           origin,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Auth:          getAuth(),
	})
}

//GetMasterBranchRefs gets the branch name
func GetMasterBranchRefs() string {
	return string(plumbing.Master)
}

//AddUpstream adds the upstream
func AddUpstream(repo *ugit.Repository, upstream string) (*ugit.Remote, error) {
	return repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{upstream},
	})
}

//FetchUpstream fetches the upstream
func FetchUpstream(remote *ugit.Remote) error {
	return remote.Fetch(&ugit.FetchOptions{
		RemoteName: "upstream",
	})
}

//PushOriginMaster push code to master
func PushOriginMaster(repo *ugit.Repository) error {
	return repo.Push(&ugit.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.DefaultPushRefSpec},
		Auth:       getAuth(),
	})
}

//CreateBranch creates branch
func CreateBranch(repo *ugit.Repository, branchName string) error {
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

//CheckoutBranch checksout branch
func CheckoutBranch(repo *ugit.Repository, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&ugit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
}

//AddCommitAndPush commits and push
func AddCommitAndPush(repo *ugit.Repository, commitMsg, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	w.Add(".")
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Rajat Jindal",
			Email: "rajatjindal83@gmail.com",
			When:  time.Now(),
		},
	})

	return repo.Push(&ugit.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.DefaultPushRefSpec},
		Auth:       getAuth(),
	})
}

//PullRebase rebases from pull
func PullRebase(repo *ugit.Repository, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       getAuth(),
	})
}
