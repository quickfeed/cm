package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func addLintCheckers(args []string) {
	fs := flag.NewFlagSet(addLintCheckersCmd, flag.ExitOnError)
	var labs string
	fs.StringVar(&labs, "labs", "", "Lab folders to add linter tests to (space separated)")

	if err := fs.Parse(args); err != nil {
		exitErr(err, "Error parsing flags")
	}
	for dir := range strings.SplitSeq(labs, " ") {
		if isLintDisabled(dir) {
			fmt.Printf("Skipping %s (lint disabled in assignment.json)\n", dir)
			continue
		}
		fmt.Printf("Updating %q in %s\n", lintFile, courseRepoPath(dir))
		if err := generateGoFromTemplate(dir, lintFile, lintTmplFile, linterTmplFS); err != nil {
			exitErr(err, "Error generating linter test")
		}
	}
}

// Assignment represents the structure of assignment.json file
type Assignment struct {
	DisableLint bool `json:"disable-lint"`
}

// isLintDisabled checks if lint is disabled for the given directory
// by reading the assignment.json file and checking the disable-lint flag.
func isLintDisabled(dir string) bool {
	assignmentPath := filepath.Join(dir, "assignment.json")
	data, err := os.ReadFile(assignmentPath)
	if err != nil {
		// If we can't read the assignment.json file, assume lint is not disabled
		return false
	}

	var assignment Assignment
	if err := json.Unmarshal(data, &assignment); err != nil {
		// If we can't parse the JSON, assume lint is not disabled
		return false
	}
	return assignment.DisableLint
}

type GoTemplateConfig struct {
	Package string
	Lab     string
}

//go:embed linter_qf_test.tmpl
var linterTmplFS embed.FS

const (
	lintTmplFile = "linter_qf_test.tmpl"
	lintFile     = "linter_qf_test.go"
)

func firstPathElem(path string) string {
	i := strings.Index(path, string(os.PathSeparator))
	if i < 0 {
		return path
	}
	return path[0:i]
}

func generateGoFromTemplate(dir, file, tmplFile string, tmplFS fs.FS) error {
	path := filepath.Join(dir, file)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFS(tmplFS, tmplFile)
	if err != nil {
		return err
	}
	config := &GoTemplateConfig{
		Package: packageName(dir),
		Lab:     firstPathElem(dir),
	}
	return tmpl.Execute(f, config)
}

// packageName returns the package name of the first Go file in the given directory.
// If no Go files are found, the base directory name is returned.
// If the directory name starts with a number, it's prefixed with "lab".
func packageName(dir string) string {
	noPkg := ensureValidPackageName(filepath.Base(dir))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return noPkg
	}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".go" {
			filePath := filepath.Join(dir, entry.Name())
			fset := token.NewFileSet()
			// parse only the package clause of the file
			f, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
			if err != nil {
				// if we can't parse the file, try the next one
				continue
			}
			return ensureValidPackageName(f.Name.Name)
		}
	}
	return noPkg
}

// ensureValidPackageName ensures the package name is valid by prefixing with "lab"
// if the name starts with a number.
func ensureValidPackageName(name string) string {
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		return "lab" + name
	}
	return name
}
