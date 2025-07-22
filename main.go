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
	Description string
	FlagUsage   string
	Function    func([]string)
}

// getCommands returns all available commands
func getCommands() map[string]Command {
	return map[string]Command{
		initEnvCmd: {
			Description: "Initialize environment variables for the course",
			FlagUsage:   "-year <year> -name <course-name> [-course <course-code>] [-discord-join-url <url>] [-bot-user <username>]",
			Function:    initEnv,
		},
		initReposCmd: {
			Description: "Initialize course repositories (assignments, info, tests)",
			FlagUsage:   "",
			Function:    initRepos,
		},
		listReposCmd: {
			Description: "List all student and group repositories",
			FlagUsage:   "[-url]",
			Function:    listRepos,
		},
		cloneRepoCmd: {
			Description: "Clone a specified repository or current user's repository",
			FlagUsage:   "[-repo <repo-name>] [-pull]",
			Function:    cloneRepo,
		},
		cloneAllReposCmd: {
			Description: "Clone all student and group repositories",
			FlagUsage:   "",
			Function:    cloneAllRepos,
		},
		updateDocTagsCmd: {
			Description: "Update documentation tags in markdown files",
			FlagUsage:   "-repo <assignments|tests|info>",
			Function:    updateDocTags,
		},
		removeSolutionTagsCmd: {
			Description: "Remove solution build tags from Go files",
			FlagUsage:   "-repo <assignments|tests>",
			Function:    removeSolutionTags,
		},
		genReadmeCmd: {
			Description: "Generate README.md files from readme_tmpl.md templates",
			FlagUsage:   "",
			Function:    genReadme,
		},
		genTestsJSONCmd: {
			Description: "Generate tests.json files for lab assignments",
			FlagUsage:   "[-view] -labs <lab1>",
			Function:    genTestsJSON,
		},
		renameLegacyTestsCmd: {
			Description: "Rename legacy *_ag_test.go files to *_qf_test.go",
			FlagUsage:   "",
			Function:    renameLegacyTests,
		},
		addLintCheckersCmd: {
			Description: "Add linter test files to specified lab folder",
			FlagUsage:   "-labs <lab1>",
			Function:    addLintCheckers,
		},
		addMainTestsCmd: {
			Description: "Add main test files to directories with existing test files",
			FlagUsage:   "",
			Function:    addMainTests,
		},
		helpCmd: {
			Description: "Show help for commands",
			FlagUsage:   "[command]",
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
		fmt.Printf("  %-*s  %s\n", maxLen, cmd, command.Description)
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
	fmt.Printf("Usage: cm %s %s\n", cmd, command.FlagUsage)
}

func exitErr(err error, msg string) {
	fmt.Printf("%s: %s\n", msg, err)
	os.Exit(1)
}
