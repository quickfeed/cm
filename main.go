package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
)

const (
	initEnvCmd            = "init-env"
	initReposCmd          = "init-repos"
	listReposCmd          = "list-repos"
	cloneRepoCmd          = "clone-repo"
	cloneAllReposCmd      = "clone-all-repos"
	updateDocTagsCmd      = "update-doc-tags"
	removeSolutionTagsCmd = "remove-solution-tags"
	genReadmeCmd          = "gen-readme"
	genTestsJSONCmd       = "gen-tests-json"
	renameLegacyTestsCmd  = "rename-tests"
	addLintCheckersCmd    = "add-lint-checkers"
	addMainTestsCmd       = "add-main-tests"
	helpCmd               = "help"
)

// Command represents a single command with its description and function
type Command struct {
	Name        string
	Description string
	Usage       string
	Function    func([]string)
}

// getCommands returns all available commands
func getCommands() map[string]Command {
	return map[string]Command{
		initEnvCmd: {
			Name:        initEnvCmd,
			Description: "Initialize environment variables for the course",
			Usage:       "cm init-env -year <year> -name <course-name> [-course <course-code>] [-discord-join-url <url>] [-bot-user <username>]",
			Function:    initEnv,
		},
		initReposCmd: {
			Name:        initReposCmd,
			Description: "Initialize course repositories (assignments, info, tests)",
			Usage:       "cm init-repos",
			Function:    func([]string) { initRepos() },
		},
		listReposCmd: {
			Name:        listReposCmd,
			Description: "List all student and group repositories",
			Usage:       "cm list-repos [-url]",
			Function:    listRepos,
		},
		cloneRepoCmd: {
			Name:        cloneRepoCmd,
			Description: "Clone a specified repository or current user's repository",
			Usage:       "cm clone-repo [-repo <repo-name>] [-pull]",
			Function:    cloneRepo,
		},
		cloneAllReposCmd: {
			Name:        cloneAllReposCmd,
			Description: "Clone all student and group repositories",
			Usage:       "cm clone-all-repos",
			Function:    func([]string) { cloneAllRepos() },
		},
		updateDocTagsCmd: {
			Name:        updateDocTagsCmd,
			Description: "Update documentation tags in markdown files",
			Usage:       "cm update-doc-tags -repo <assignments|tests|info>",
			Function:    updateDocTags,
		},
		removeSolutionTagsCmd: {
			Name:        removeSolutionTagsCmd,
			Description: "Remove solution build tags from Go files",
			Usage:       "cm remove-solution-tags -repo <assignments|tests>",
			Function:    removeSolutionTags,
		},
		genReadmeCmd: {
			Name:        genReadmeCmd,
			Description: "Generate README.md files from readme_tmpl.md templates",
			Usage:       "cm gen-readme",
			Function:    func([]string) { genReadme() },
		},
		genTestsJSONCmd: {
			Name:        genTestsJSONCmd,
			Description: "Generate tests.json files for lab assignments",
			Usage:       "cm gen-tests-json [-view] -labs <lab1>",
			Function:    genTestsJSON,
		},
		renameLegacyTestsCmd: {
			Name:        renameLegacyTestsCmd,
			Description: "Rename legacy *_ag_test.go files to *_qf_test.go",
			Usage:       "cm rename-tests",
			Function:    func([]string) { renameLegacyTests() },
		},
		addLintCheckersCmd: {
			Name:        addLintCheckersCmd,
			Description: "Add linter test files to specified lab folder",
			Usage:       "cm add-lint-checkers -labs <lab1>",
			Function:    addLintCheckers,
		},
		addMainTestsCmd: {
			Name:        addMainTestsCmd,
			Description: "Add main test files to directories with existing test files",
			Usage:       "cm add-main-tests",
			Function:    func([]string) { addMainTests() },
		},
		helpCmd: {
			Name:        helpCmd,
			Description: "Show help for commands",
			Usage:       "cm help [command]",
			Function:    showHelp,
		},
	}
}

func main() {
	if len(os.Args) < 2 {
		usageMsg()
		return
	}

	cmd, args := os.Args[1], os.Args[2:]
	commands := getCommands()
	if command, exists := commands[cmd]; exists {
		command.Function(args)
		return
	}
	fmt.Printf("Unknown command: %s\n", cmd)
	usageMsg()
}

func usageMsg() {
	fmt.Println("Usage: cm <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")

	commands := getCommands()
	commandNames := slices.Sorted(maps.Keys(commands))
	// find the longest command name for formatting
	maxLen := len(slices.MaxFunc(commandNames, func(a, b string) int {
		return len(a) - len(b)
	}))

	// print commands with descriptions in sorted order
	for _, cmd := range commandNames {
		command := commands[cmd]
		fmt.Printf("  %-*s  %s\n", maxLen, command.Name, command.Description)
	}

	fmt.Println()
	fmt.Println("Use 'cm help <command>' for detailed help on a specific command.")
	os.Exit(1)
}

func showHelp(args []string) {
	if len(args) == 0 {
		usageMsg()
		return
	}

	cmd := args[0]
	commands := getCommands()
	command, exists := commands[cmd]
	if !exists {
		fmt.Printf("Unknown command: %s\n", cmd)
		usageMsg()
		return
	}
	fmt.Printf("Description: %s\n", command.Description)
	fmt.Printf("Usage: %s\n", command.Usage)
}

func exitErr(err error, msg string) {
	fmt.Printf("%s: %s\n", msg, err)
	os.Exit(1)
}
