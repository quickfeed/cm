package main

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	courseRepos     = []string{"assignments", "info", "tests"}
	assignmentsRepo = courseRepos[0]
)

// initRepos initializes the course repositories so that they are
// ready to be populated with data from the main course repository.
// Once the repositories are initialized, they are ready to be pushed
// to the per-year course organization on GitHub.
func initRepos(_ []string) {
	if err := loadEnv(); err != nil {
		exitErr(err, "Error loading environment variables")
	}

	for _, repo := range courseRepos {
		path := repoPath(repo)
		if exists(path) {
			fmt.Printf("Repository %q already exists, skipping.\n", repo)
			continue
		}
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			exitErr(err, "Error creating repository directory")
		}
		if err := initGitRepository(path, courseOrg(), repo); err != nil {
			exitErr(err, "Error initializing git repository")
		}
	}
}

func initGitRepository(workingDir, ghOrg, repo string) error {
	commands := [][]string{
		{"git", "init"},
		{"git", "branch", "-M", "main"},
		{"git", "remote", "add", "origin", gitURL(ghOrg, repo)},
	}

	// Try to pull if remote has content
	if remoteHasContent(workingDir, "origin", "main") {
		commands = append(commands,
			[]string{"git", "pull", "origin", "main"},
			[]string{"git", "branch", "--set-upstream-to=origin/main", "main"},
		)
	}

	for _, cmd := range commands {
		if err := runCommand(workingDir, cmd...); err != nil {
			return err
		}
	}
	return nil
}

func remoteHasContent(workingDir, remote, branch string) bool {
	cmd := exec.Command("git", "ls-remote", "--heads", remote, branch)
	cmd.Dir = workingDir
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

func gitURL(ghOrg, repo string) string {
	// TODO(meling) consider supporting HTTPS remote origin
	return fmt.Sprintf("git@github.com:%s/%s.git", ghOrg, repo)
}
