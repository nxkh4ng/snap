package utils

import (
	"fmt"
	"os/exec"
)

func CheckGitRepo() error {
	repoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := repoCmd.Run(); err != nil {
		return fmt.Errorf("not a git repository, use `git init` to initialize")
	}
	return nil
}

func CheckGitCommitReady() error {
	stagedCmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := stagedCmd.Run()
	if err == nil {
		return fmt.Errorf("no staged changes to commit, use `git add <file>` to stage your changes")
	} else {
		if err.Error() != "exit status 1" {
			return fmt.Errorf("failed to check staged changes: %v", err)
		}
	}
	return nil
}

func StageAll() error {
	addCmd := exec.Command("git", "add", "-u")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage files: %v", err)
	}
	return nil
}

func GetTheLatestCommitMsg() (string, error) {
	logCmd := exec.Command("git", "log", "-1", "--format=%B")
	out, err := logCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to log latest commit: %v", err)
	}
	return string(out), nil
}

func Commit(msg string, amendFlag bool) (string, error) {
	var cmd *exec.Cmd
	if amendFlag {
		cmd = exec.Command("git", "commit", "--amend", "-m", msg)
	} else {
		cmd = exec.Command("git", "commit", "-m", msg)
	}

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to commit: %v", err)
	}
	return string(out), nil
}
