package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// GetInfo retorna a branch atual e o hash curto do commit
func GetInfo() (string, string, string) {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	branch := strings.TrimSpace(string(branchOut))
	if err != nil {
		branch = "unknown"
	}

	commitCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	commitOut, err := commitCmd.Output()
	commit := strings.TrimSpace(string(commitOut))
	if err != nil {
		commit = "unknown"
	}

	topCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	topOut, err := topCmd.Output()
	project := strings.TrimSpace(string(topOut))
	if err != nil || project == "" {
		project = "unknown"
	} else {
		project = filepath.Base(project)
	}

	return branch, commit, project
}
