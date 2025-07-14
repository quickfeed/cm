package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

// runCommandWithOutput runs a command and parses the output into a specified type.
// It is useful for commands that return JSON output.
//
// Example usage:
//
//	type MyOutput struct {
//	    Field1 string `json:"field1"`
//	    Field2 int    `json:"field2"`
//	}
//
//	output, err := runCommandWithOutput[MyOutput]("mycommand", "arg1", "arg2"); err != nil {
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
