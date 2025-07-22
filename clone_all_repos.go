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
func cloneAllRepos(_ []string) {
	if err := loadEnv(); err != nil {
		exitErr(err, "Error loading environment variables")
	}
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
		go func(repo repositoryInfo) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if exists(repoPath(repo.Name)) {
				// If the repository already exists, we will update it instead of cloning
				if err := runContextCommand(ctx, repoPath(repo.Name), "git", "pull"); err != nil {
					fmt.Printf("Error updating %q: %v\n", repo, err)
					return
				}
				fmt.Printf("Updated %q in %q\n", repo.Name, repoPath(repo.Name))
				return
			}

			// clone the repository
			err := runContextCommand(ctx, ".", "git", "clone", repo.URL, repoPath(repo.Name))
			if err != nil {
				fmt.Printf("Error cloning %q: %v\n", repo, err)
				return
			}
			fmt.Printf("Cloned %q into %q\n", repo.Name, repoPath(repo.Name))
		}(repo)
	}

	// Wait for all clones to complete
	wg.Wait()
}
