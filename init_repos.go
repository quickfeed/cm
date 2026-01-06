package main

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	courseRepos     = []string{"assignments", "website", "tests"}
	assignmentsRepo = courseRepos[0]
)

// repoName returns the actual repository name for the given internal name.
// This allows mapping internal names (like "website") to actual repository
// names (like "dat520.github.io").
func repoName(internalName string) string {
	if internalName == "website" {
		return websiteRepo()
	}
	return internalName
}

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
		actualRepoName := repoName(repo)
		if err := initGitRepository(path, courseOrg(), actualRepoName); err != nil {
			exitErr(err, "Error initializing git repository")
		}
	}
}

func initGitRepository(workingDir, ghOrg, repo string) error {
	remoteURL := gitURL(ghOrg, repo)

	// Initialize repository and add remote
	initCommands := [][]string{
		{"git", "init"},
		{"git", "branch", "-M", "main"},
		{"git", "remote", "add", "origin", remoteURL},
	}
	for _, cmd := range initCommands {
		if err := runCommand(workingDir, cmd...); err != nil {
			return err
		}
	}

	// Try to pull if remote has content
	if remoteHasContent(remoteURL, "main") {
		pullCommands := [][]string{
			{"git", "pull", "origin", "main"},
			{"git", "branch", "--set-upstream-to=origin/main", "main"},
		}
		for _, cmd := range pullCommands {
			if err := runCommand(workingDir, cmd...); err != nil {
				return err
			}
		}
	}

	return nil
}

func remoteHasContent(remoteURL, branch string) bool {
	cmd := exec.Command("git", "ls-remote", "--heads", remoteURL, branch)
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

func gitURL(ghOrg, repo string) string {
	// TODO(meling) consider supporting HTTPS remote origin
	return fmt.Sprintf("git@github.com:%s/%s.git", ghOrg, repo)
}
