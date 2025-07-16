package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

// runCommandWithOutput runs a command and unmarshal the output into a specified type.
// It is useful for commands that return JSON output.
//
// Example usage:
//
//	type MyOutput struct {
//	    Field1 string `json:"field1"`
//	    Field2 int    `json:"field2"`
//	}
//
//	output, err := runCommandWithOutput[MyOutput](workDir, "mycommand", "arg1", "arg2"); err != nil {
//	    fmt.Println("Error:", err)
//	}
//
// Example output from mycommand might be:
//
//	{"field1": "value1", "field2": 42}
//
// The output will be parsed into the MyOutput struct.
// Example return value:
//
//	&MyOutput{Field1: "value1", Field2: 42}
func runCommandWithOutput[T any](dir string, args ...string) (result T, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("command %q failed: %w", args, err)
	}
	if !json.Valid(out) {
		// if the output is not valid JSON, check if individual lines are
		// valid JSON; if so, create a JSON array from the combined lines.
		// this is useful for commands that output multiple JSON objects, one per line.
		jsonLines := []string{"["}
		for line := range bytes.Lines(out) {
			if len(line) == 0 || !json.Valid(line) {
				continue // skip empty lines or lines that are not valid JSON
			}
			jsonLines = append(jsonLines, string(line)+",")
		}
		lastLine := len(jsonLines) - 1
		jsonLines[lastLine] = strings.TrimSuffix(jsonLines[lastLine], ",") // remove trailing comma from last line
		jsonLines = append(jsonLines, "]")
		out = []byte(strings.Join(jsonLines, "\n"))
	}

	// unmarshal the output into the specified type
	if err := json.Unmarshal(out, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	return result, nil
}
