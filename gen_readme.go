package main

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

const (
	readmeTmplFile = "readme_tmpl.md"
	readmeFile     = "README.md"
	assignmentFile = "assignment.json"
)

// LabHeader contains a string representation of the content to be written as a lab header
type LabHeader struct {
	LabHeader string
	ToC       []string
}

func genReadme(_ []string) {
	if err := loadEnv(); err != nil {
		exitErr(err, "Error loading environment variables")
	}

	labs, err := findLabsWithReadmeTmpl(gitRoot)
	if err != nil {
		exitErr(err, "Error finding labs with readme_tmpl.md files")
	}
	assignments := generateReadme(gitRoot, labs)

	labPlan := mustExecute(parseTemplate("labplan", labPlanTemplate), assignments)
	err = os.WriteFile(filepath.Join(gitRoot, "website", "lab-plan.md"), []byte(labPlan), 0o644)
	if err != nil {
		exitErr(err, "Error writing lab-plan.md")
	}
}

func generateReadme(repo string, labs map[string][]string) map[int]*AssignmentInfo {
	// process each assignment and corresponding readme_tmpl.md
	assignments := make(map[int]*AssignmentInfo)
	for _, lab := range slices.Sorted(maps.Keys(labs)) {
		for _, readme := range labs[lab] {
			// prepare paths for log output and for output file
			assignPath := strings.Replace(filepath.Join(lab, assignmentFile), repo, course(), 1)
			readmePath := strings.Replace(readme, repo, course(), 1)
			readmeMDPath := strings.Replace(readme, readmeTmplFile, readmeFile, 1)
			readmeMDPathShort := strings.Replace(readmeMDPath, repo, course(), 1)
			fmt.Printf("Combining (%v, %v) -> %v\n", assignPath, readmePath, readmeMDPathShort)

			header := ""
			if filepath.Dir(assignPath) == filepath.Dir(readmePath) {
				// Only include header for labX/readme_tmpl.md files, not those in subdirectories
				header = parseAssignmentHeader(lab, headerTemplate, assignments)
			} else {
				// Only include title for readme_tmpl.md files in subdirectories
				header = parseAssignmentHeader(lab, titleOnlyHeaderTemplate, assignments)
			}

			// load readme_tmpl.md and execute template; saving output to README.md
			b, err := os.ReadFile(readme)
			check(err)
			readme = string(b)
			// compute position in readme without the level one title, if exists
			if !strings.HasPrefix(readme, "# ") {
				panic(readmePath + " expected to start with markdown title (single #)")
			}
			posWithoutTitle := strings.Index(readme, "\n")
			if posWithoutTitle > 0 && readme[posWithoutTitle] == '\n' {
				posWithoutTitle++
			}
			// inject the template lines in the beginning of the readme
			readmeTemplate := fmt.Sprintf("%v\n%v", tocTemplate, readme[posWithoutTitle+1:])
			toc := generateToC(readmeTemplate)
			readmeHeader := &LabHeader{
				LabHeader: header,
				ToC:       toc,
			}
			tocMD := mustExecute(parseTemplate("readme", tocTemplate), readmeHeader)
			readmeMD := tocMD + readme[posWithoutTitle+1:]
			err = os.WriteFile(readmeMDPath, []byte(readmeMD), 0o644)
			check(err)
		}
	}
	return assignments
}

// findLabsWithReadmeTmpl returns a map of labs with assignment.json files
// and a slice of their corresponding readme_tmpl.md files.
func findLabsWithReadmeTmpl(repo string) (map[string][]string, error) {
	labs := make(map[string][]string)

	// find all labs with assignment.json files and check for legacy files
	err := filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		var emptySlice []string
		if d.Name() == assignmentFile {
			dir := filepath.Dir(path)
			labs[dir] = emptySlice
		} else if d.Name() == "assignment.yml" || d.Name() == "assignment.yaml" {
			fmt.Printf("Warning: Found legacy '%s' file. Run 'cm convert-yaml-to-json' to convert.\n", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// find all corresponding readme_tmpl.md files
	err = filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == readmeTmplFile {
			dir := filepath.Dir(path)
			if _, found := labs[dir]; found {
				labs[dir] = append(labs[dir], path)
			} else {
				for level := 4; !found && level > 0; level-- {
					// traverse up the hierarchy looking for existing lab dir
					// with a previously recorded assignment.json file;
					// stop when level reach 0
					dir := filepath.Dir(dir)
					if _, found = labs[dir]; found {
						labs[dir] = append(labs[dir], path)
					}
				}
				if !found {
					fmt.Printf("ignoring %s: couldn't find %s in or above %v\n", path, assignmentFile, dir)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return labs, nil
}

// parseAssignmentHeader returns a header by parsing assignment.json.
func parseAssignmentHeader(lab, headerTemplate string, assignments map[int]*AssignmentInfo) string {
	assignment, err := parseAssignment(filepath.Join(lab, assignmentFile))
	check(err)

	// make sure all assignments has CourseOrg field set
	assignment.CourseOrg = course()
	if _, found := assignments[assignment.Order]; !found {
		// add to assignments only once; this is since the assignment.json
		// may exist for multiple versions of the same assignment README.md.
		assignments[assignment.Order] = assignment
	}
	return mustExecute(parseTemplate("assignment", headerTemplate), assignment)
}

// stripMarkdownLinks removes markdown link syntax [text](url) from a string,
// keeping only the text content.
func stripMarkdownLinks(text string) string {
	// Remove markdown links [text](url) and keep only the text
	linkRegExp := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	return linkRegExp.ReplaceAllString(text, "$1")
}

// generateToC takes a markdown file and generates a table of contents
// of all level two and three headings, preserving the heading level.
func generateToC(readme string) []string {
	// this reg exp became a bit nasty since we want to match with backtick
	// now also includes parentheses and square brackets for markdown links
	legalHeadingChars := `\w\s\/:#-\[\]()` + "`"
	headingRegExp := regexp.MustCompile(`^(#{2,3})\s([` + legalHeadingChars + `]+)$`)
	headings := make([]string, 0)
	for line := range strings.SplitSeq(readme, "\n") {
		if headingRegExp.MatchString(line) {
			// Strip markdown links from the heading before adding to ToC
			cleanedLine := stripMarkdownLinks(line)
			headings = append(headings, cleanedLine)
		}
	}
	return headings
}

func headerLabels(data *AssignmentInfo) []string {
	labels := []string{
		fmt.Sprintf("Lab %d:", data.Order),
		"Subject:",
		"Deadline:",
	}
	if data.Effort != "" {
		labels = append(labels, "Expected effort:")
	}
	if data.ScoreLimit > 0 {
		labels = append(labels, "Score limit:")
	}
	return append(labels, "Grading:", "Submission:")
}

func headerValues(data *AssignmentInfo) []string {
	values := []string{
		data.Title,
		data.Subject,
		fmt.Sprintf("**%s**", data.Deadline),
	}
	if data.Effort != "" {
		values = append(values, data.Effort)
	}
	if data.ScoreLimit > 0 {
		values = append(values, fmt.Sprintf("%d", data.ScoreLimit))
	}
	return append(values, data.Grading, data.SubmissionType)
}

var funcMap = template.FuncMap{
	"link": func(heading string) string {
		// First strip any remaining markdown links
		linkRegExp := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
		heading = linkRegExp.ReplaceAllString(heading, "$1")

		replace := map[string]string{
			" ":  "-",
			"#":  "",
			"`":  "",
			":":  "",
			"/":  "",
			".":  "",
			",":  "",
			"(":  "",
			")":  "",
			"[":  "",
			"]":  "",
			"&":  "",
			"**": "",
			"+":  "",
			"?":  "",
			"=":  "",
		}
		str := strings.ToLower(heading)
		for old, new := range replace {
			str = strings.ReplaceAll(str, old, new)
		}
		// Clean up double dashes that can occur when & is between spaces
		str = strings.ReplaceAll(str, "--", "-")
		return str
	},
	"escapeText": func(heading string) string {
		// First strip any markdown links
		linkRegExp := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
		heading = linkRegExp.ReplaceAllString(heading, "$1")
		// Escape & to \& for markdown link text
		return strings.ReplaceAll(heading, "&", "\\&")
	},
	"inc": func(i int) int {
		return i + 1
	},
	"add": func(a, b int) int {
		return a + b
	},
	"hasPrefix": func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	},
	"trimPrefix": func(s, prefix string) string {
		return strings.TrimPrefix(s, prefix)
	},
	"headerLabels": headerLabels,
	"headerValues": headerValues,
	"col1Width": func(data *AssignmentInfo) int {
		return len(slices.MaxFunc(headerLabels(data), func(a, b string) int {
			return len(a) - len(b)
		}))
	},
	"col2Width": func(data *AssignmentInfo) int {
		return len(slices.MaxFunc(headerValues(data), func(a, b string) int {
			return len(a) - len(b)
		}))
	},
	"padRight": func(s string, width int) string {
		if len(s) >= width {
			return s
		}
		return s + strings.Repeat(" ", width-len(s))
	},
	"center": func(s string, width int) string {
		if len(s) >= width {
			return s
		}
		padding := width - len(s)
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
	},
	"separator": func(width int) string {
		return strings.Repeat("-", width)
	},
}

func parseTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Funcs(funcMap).Parse(tmpl))
}

func mustExecute(t *template.Template, data any) string {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		panic(err)
	}
	return b.String()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
