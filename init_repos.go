package main

import (
	"context"
	"encoding/json"
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
func initRepos() {
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
	for _, cmd := range commands {
		if err := runCommand(workingDir, cmd...); err != nil {
			return err
		}
	}
	return nil
}

func runCommand(dir string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w", args, err)
	}
	return nil
}

func runContextCommand(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w", args, err)
	}
	return nil
}

func runCommandWithOutput[T any](args ...string) (*T, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command %q failed: %w", args, err)
	}
	// parse the output if needed
	// for example, if the command returns JSON, you can unmarshal it into a struct
	var result T
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	return &result, nil
}

func gitURL(ghOrg, repo string) string {
	// TODO(meling) consider supporting HTTPS remote origin
	return fmt.Sprintf("git@github.com:%s/%s.git", ghOrg, repo)
}
