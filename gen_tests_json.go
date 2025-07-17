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
	goTestCommand = "go test -v -tags solution ./..."
)

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

func genTestsJSON(args []string) {
	fs := flag.NewFlagSet(genTestsJSONCmd, flag.ExitOnError)
	labs := fs.String("labs", "", "Lab folders to generate tests.json for (space separated)")
	view := fs.Bool("view", false, "Show the JSON output of the generated tests.json file")

	if err := fs.Parse(args); err != nil {
		exitErr(err, "Error parsing flags")
	}

	for dir := range strings.SplitSeq(*labs, " ") {
		if !exists(courseRepoPath(dir)) {
			exitErr(fmt.Errorf("directory %q does not exist", dir), errMsg)
		}
		// Annoyingly, we need to run this with both -v and ./... to get the JSON output.
		// The ./... because some labs contain several subdirectories with tests.
		// Running without the -v flag, it will suppress the output and we won't get the JSON.
		scoreList, err := runCommandWithOutput[[]score](courseRepoPath(dir), strings.Split(goTestCommand, " ")...)
		if scoreList == nil {
			exitErr(fmt.Errorf("no score objects found for %q", dir), errMsg)
		}
		if err != nil {
			fmt.Printf("Error running %q for %s: %v\n", goTestCommand, dir, err)
			fmt.Printf("You may want to run %q manually to debug the issue.\n", goTestCommand)
			// Try to process the scores anyway; this will allow us to generate a tests.json
			// file even if the command fails. This is usually fine since the empty score
			// objects that we need are usually printed at the start of the test output.
			// The exception could be if there was a compile error or panic situation.
		}

		scoreList = deduplicateScores(scoreList)

		// Write the scores to tests.json
		file := filepath.Join(courseRepoPath(dir), testsJSONFile)
		f, err := os.Create(file)
		if err != nil {
			exitErr(err, errMsg)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(scoreList); err != nil {
			exitErr(err, errMsg)
		}
		fmt.Printf("Generated %s for %s\n", testsJSONFile, courseRepoPath(dir))
		if *view {
			for _, s := range scoreList {
				fmt.Printf("{%q:%q,%q:%d,%q:%d}\n", "TestName", s.TestName, "Max Score", s.MaxScore, "Weight", s.Weight)
			}
		}
	}
}
