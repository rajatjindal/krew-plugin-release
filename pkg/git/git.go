package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	ugit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

//Clone clones the repo
func Clone(origin, branch, dir string) error {
	auth := &githttp.BasicAuth{
		Username: "rjindal",
		Password: os.Getenv("KREW_PLUGIN_RELEASE_TOKEN"),
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

func GetMasterBranchRefs() string {
	return string(plumbing.Master)
}

func AddUpstream(dir, upstream string) error {
	cmdline := []string{
		"remote",
		"add",
		"upstream",
		upstream,
	}

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git remote add upstream command.")
	}

	return nil
}

func FetchUpstream(dir string) error {
	cmdline := []string{
		"fetch",
		"upstream",
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git fetch upstream command.")
	}

	return nil
}

func RebaseUpstream(dir string) error {
	cmdline := []string{
		"rebase",
		"upstream/master",
	}

	fmt.Println(dir, cmdline)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git rebase upstream master command.")
	}

	return nil
}

func PushOriginMaster(dir string) error {
	cmdline := []string{
		"push",
		"origin",
		"master",
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git push origin master command.")
	}

	return nil
}

func CreateBranch(dir, branchName string) error {
	cmdline := []string{
		"branch",
		branchName,
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git branch command.")
	}

	return nil
}

func CheckoutBranch(dir, branchName string) error {
	cmdline := []string{
		"checkout",
		branchName,
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git checkout command.")
	}

	return nil
}

func CommitAndPush(dir, commitMsg, branchName string) error {
	err := Commit(dir, commitMsg)
	if err != nil {
		return err
	}

	err = Push(dir, branchName)
	if err != nil {
		return err
	}

	return nil
}

func Commit(dir, commitMsg string) error {
	cmdline := []string{
		"commit",
		"-m",
		commitMsg,
		".",
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git commit command.")
	}

	return nil
}

func Push(dir, branchName string) error {
	cmdline := []string{
		"push",
		"origin",
		branchName,
	}

	fmt.Println(dir)

	cmd := exec.Command("git", cmdline...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed when running git push command.")
	}

	return nil
}

func ForkRepo(upstream string) error {
	return nil
}
