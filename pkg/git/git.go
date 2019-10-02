package git

import (
	"io/ioutil"
	"os"

	ssh "golang.org/x/crypto/ssh"
	ugit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

//GetAuthKeys returns the auth keys used for git ops
func GetAuthKeys(key string) (*gitssh.PublicKeys, error) {
	keyContent, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}

	signer, _ := ssh.ParsePrivateKey([]byte(keyContent))
	return &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
	}, nil
}

//Clone clones the repo
func Clone(origin, branch, dir string) error {
	_, err := ugit.PlainClone(dir, false, &ugit.CloneOptions{
		URL:           origin,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
	})

	if err != nil {
		return err
	}

	return nil
}
