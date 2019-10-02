package git

import (
	"fmt"
	"net"
	"os"

	ssh "golang.org/x/crypto/ssh"
	ugit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

//Clone clones the repo
func Clone(origin, branch, dir string) error {
	if len(os.Getenv("GITHUB_SSH_KEY")) == 0 {
		return fmt.Errorf("env GITHUB_SSH_KEY not found")
	}

	signer, _ := ssh.ParsePrivateKey([]byte(os.Getenv("GITHUB_SSH_KEY")))
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				fmt.Printf("verifying for hostname %s, remote %s\n", hostname, remote.String())
				return nil
			},
		},
	}

	_, err := ugit.PlainClone(dir, false, &ugit.CloneOptions{
		URL:           origin,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Auth:          auth,
	})

	if err != nil {
		return err
	}

	return nil
}

func CreateBranch(dir, branchName string) error {
	return nil
}

func CommitAndPush(dir, origin string) error {
	return nil
}
