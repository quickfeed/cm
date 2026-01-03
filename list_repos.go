package main

import (
	"flag"
	"fmt"
)

type repositoryInfo struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// listRepos lists all student and group repositories from a GitHub course.
//
// Run with the following command:
//
//	cm list-repos
//	cm list-repos -url
//
// The -url flag specifies that the URL of the repository should be printed instead of the name.
func listRepos(args []string) {
	if err := loadEnv(); err != nil {
		exitErr(err, "Error loading environment variables")
	}
	fs := flag.NewFlagSet(cloneRepoCmd, flag.ExitOnError)
	var url bool
	fs.BoolVar(&url, "url", false, "Print only the URL of the repository")

	if err := fs.Parse(args); err != nil {
		exitErr(err, "Error parsing flags")
	}

	repos, err := getRepositories(courseOrg())
	if err != nil {
		exitErr(err, "Error listing repositories")
	}
	for _, repo := range repos {
		fmt.Printf("%s", repo.Name)
		if url {
			fmt.Printf(" (%s)", repo.URL)
		}
		fmt.Println()
	}
}

// getRepositories returns a list of repositories for the given organization.
// It uses the GitHub CLI command to fetch the repositories with the following command:
//
//	gh repo list <org> --limit 100 --json name,url
func getRepositories(org string) ([]repositoryInfo, error) {
	repos, err := runCommandWithOutput[[]repositoryInfo](gitRoot, "gh", "repo", "list", org, "--limit", "100", "--json", "name,url")
	if err != nil {
		exitErr(err, "Error fetching repositories")
	}
	if len(repos) == 0 {
		exitErr(fmt.Errorf("no repositories found"), "Error listing repositories")
	}
	return repos, nil
}
