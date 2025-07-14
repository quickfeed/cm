package main

import (
	"flag"
	"fmt"
)

type listReposJSON struct {
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

	// Alternative: use the GitHub CLI to list repositories:
	// gh repo list dat520-2025 --limit 100 --json name,url
	// gh repo list dat520-2025 --limit 100 --json name
	for _, repo := range repos {
		fmt.Printf("%s", repo.Name)
		if url {
			fmt.Printf(" (%s)", repo.URL)
		}
		fmt.Println()
	}
}

func getRepositories(org string) ([]listReposJSON, error) {
	repos, err := runCommandWithOutput[[]listReposJSON]("gh", "repo", "list", org, "--limit", "100", "--json", "name,url")
	if err != nil {
		exitErr(err, "Error fetching repositories")
	}
	if repos == nil || len(*repos) == 0 {
		exitErr(fmt.Errorf("no repositories found"), "Error listing repositories")
	}
	return *repos, nil
}
