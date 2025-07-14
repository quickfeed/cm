package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

// cloneAllRepos clones all student and group repositories from a GitHub course.
//
// Run with the following command:
//
//	cm clone-all-repos
//
// The student/group repositories for the current year are cloned into the directory
// <year>/student-repos/<repo>.
func cloneAllRepos() {
	if err := loadEnv(); err != nil {
		exitErr(err, "Error loading environment variables")
	}
	// Alternative: use the GitHub CLI to list repositories:
	// gh repo list dat520-2025 --limit 100 --json name,url
	// gh repo list dat520-2025 --limit 100 --json name
	ghRepos, err := getRepositories(courseOrg())
	if err != nil {
		exitErr(err, "Error listing repositories")
	}

	path := filepath.Join(repoHome(), "student-repos")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		exitErr(err, "Error creating repository directory")
	}
	repoPath := func(repo string) string {
		return filepath.Join(path, repo)
	}

	fmt.Printf("Cloning %d student/group repositories into %q\n", len(ghRepos), path)

	// Create a semaphore to limit concurrent clones to 5
	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for _, repo := range ghRepos {
		// skipping the main course repository (assignments, info, tests)
		if slices.Contains(courseRepos, repo.Name) {
			continue
		}

		wg.Add(1)
		go func(repo listReposJSON) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore

			msg := fmt.Sprintf("Cloned %q into %q", repo, path)
			if exists(repoPath(repo.Name)) {
				msg = fmt.Sprintf("Repository %q updated", repo)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// clone the repository
			err := runContextCommand(ctx, ".", "git", "clone", repo.URL, repoPath(repo.Name))
			if err != nil {
				fmt.Printf("Error cloning %q: %v\n", repo, err)
				return
			}
			fmt.Println(msg)
		}(repo)
	}

	// Wait for all clones to complete
	wg.Wait()
}
