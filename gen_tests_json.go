package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	testsJSONFile = "tests.json"
	errMsg        = "No tests.json file generated"
)

// testCommands holds the test commands for supported languages.
// The language key should be lowercase.
var testCommands = map[string]string{
	"go": "go test -v -tags solution ./...",
	"c#": "dotnet -v",
}

type score struct {
	TestName string
	MaxScore int
	Weight   int
}

// deduplicateScores returns a slice of unique score objects based on TestName.
func deduplicateScores(scores []score) []score {
	seenMap := make(map[string]score)
	for _, s := range scores {
		if _, exists := seenMap[s.TestName]; !exists {
			seenMap[s.TestName] = s
		}
	}
	uniqueScores := slices.Collect(maps.Values(seenMap))
	slices.SortFunc(uniqueScores, func(a, b score) int {
		return strings.Compare(a.TestName, b.TestName)
	})
	return uniqueScores
}

// formatCompactJSON formats the score slice as JSON with each test entry on a single line
func formatCompactJSON(scores []score) ([]byte, error) {
	var lines []string
	lines = append(lines, "[")

	for i, s := range scores {
		// Marshal each score object to a single line
		data, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}

		line := "  " + string(data)
		if i < len(scores)-1 {
			line += ","
		}
		lines = append(lines, line)
	}

	lines = append(lines, "]")
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}

func genTestsJSON(args []string) {
	fs := flag.NewFlagSet(genTestsJSONCmd, flag.ExitOnError)
	labs := fs.String("labs", "", "Lab folders to generate tests.json for (space separated)")
	view := fs.Bool("view", false, "Show the JSON output of the generated tests.json file")
	runnerLang := fs.String("runner", "go", "Test runner's programming language (default: go)")

	if err := fs.Parse(args); err != nil {
		exitErr(err, "Error parsing flags")
	}

	lang := strings.ToLower(*runnerLang)
	testRunnerCmd, ok := testCommands[lang]
	if !ok {
		languages := slices.Collect(maps.Keys(testCommands))
		exitErr(fmt.Errorf("unsupported language %q, supported languages are: %v", lang, languages), errMsg)
	}

	for dir := range strings.SplitSeq(*labs, " ") {
		if !exists(courseRepoPath(dir)) {
			exitErr(fmt.Errorf("directory %q does not exist", dir), errMsg)
		}
		// We can should run the test command with the SCORE_INIT environment variable set to 1,
		// which will avoid running the tests and only initialize the score objects.
		if err := os.Setenv("SCORE_INIT", "1"); err != nil {
			exitErr(err, errMsg)
		}
		// Annoyingly, we need to run this with both -v and ./... to get the JSON output.
		// The ./... because some labs contain several subdirectories with tests.
		// Running without the -v flag, it will suppress the output and we won't get the JSON.
		scoreList, err := runCommandWithOutput[[]score](courseRepoPath(dir), strings.Split(testRunnerCmd, " ")...)
		if err != nil {
			fmt.Printf("Error running %q for %s: %v\n", testRunnerCmd, dir, err)
			fmt.Printf("You may want to run %q manually to debug the issue.\n", testRunnerCmd)
			// Try to process the scores anyway; this will allow us to generate a tests.json
			// file even if the command fails. This is usually fine since the empty score
			// objects that we need are usually printed at the start of the test output.
			// The exception could be if there was a compile error or panic situation.
		}

		// Skip generating tests.json if there are no tests
		if len(scoreList) == 0 {
			fmt.Printf("No tests found for %s, skipping %s generation\n", dir, testsJSONFile)
			continue
		}

		scoreList = deduplicateScores(scoreList)

		// Write the scores to tests.json
		file := filepath.Join(courseRepoPath(dir), testsJSONFile)
		f, err := os.Create(file)
		if err != nil {
			exitErr(err, errMsg)
		}
		defer f.Close()

		// Use custom formatter to keep each test entry on a single line
		jsonData, err := formatCompactJSON(scoreList)
		if err != nil {
			exitErr(err, errMsg)
		}
		if _, err := f.Write(jsonData); err != nil {
			exitErr(err, errMsg)
		}
		fmt.Printf("Generated %s for %s\n", testsJSONFile, courseRepoPath(dir))
		if *view {
			for _, s := range scoreList {
				jsonData, err := json.Marshal(s)
				if err != nil {
					continue
				}
				fmt.Println(string(jsonData))
			}
		}
	}
}
