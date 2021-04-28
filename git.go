package gobuild

import (
	"os/exec"
	"strings"
)

// GetGitVersion returns the current git tag,
// which is the result of `git describe --tags` run in `path`
func GetGitVersion(path string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitHash returns the current git commit hash,
// which is the result of `git rev-parse HEAD` run in `path`
func GetGitHash(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
