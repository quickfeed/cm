package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// convertYamlToJSON converts assignment.yml or assignment.yaml files to JSON format.
// It does this by simply wrapping the YAML content with { } and validating the result is valid JSON.
// If successful, it writes the content to assignment.json and deletes the original YAML file.
func convertYamlToJSON(args []string) {
	yamlFiles, err := findAssignmentYamlFiles(gitRoot)
	if err != nil {
		exitErr(err, "Error finding YAML files")
	}
	if len(yamlFiles) == 0 {
		fmt.Println("No assignment.yml or assignment.yaml files found.")
		return
	}

	for _, yamlFile := range yamlFiles {
		if err := convertSingleFile(yamlFile); err != nil {
			fmt.Printf("Error converting %s: %s\n", yamlFile, err)
			continue
		}
		fmt.Printf("Successfully converted %s to assignment.json\n", yamlFile)
	}
}

// findAssignmentYamlFiles finds all assignment.yml and assignment.yaml files in the given directory and its subdirectories
func findAssignmentYamlFiles(path string) ([]string, error) {
	var yamlFiles []string
	err := filepath.WalkDir(path, func(filePath string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dirEntry.IsDir() {
			return nil
		}
		filename := dirEntry.Name()
		if filename == "assignment.yml" || filename == "assignment.yaml" {
			yamlFiles = append(yamlFiles, filePath)
		}
		return nil
	})
	return yamlFiles, err
}

// convertSingleFile converts a YAML file to JSON.
func convertSingleFile(yamlFilePath string) error {
	content, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}
	contentStr := strings.TrimSpace(string(content))
	if contentStr == "" {
		return fmt.Errorf("empty YAML file")
	}

	// Convert YAML format to JSON format by adding quotes around keys
	jsonContent := convertYamlKeysToJson(contentStr)
	
	// Wrap content with { } to make it JSON-like
	jsonContent = "{" + jsonContent + "}"

	// Validate that the result is valid JSON and marshal it back to ensure formatting
	var jsonObj any
	if err := json.Unmarshal([]byte(jsonContent), &jsonObj); err != nil {
		return fmt.Errorf("resulting content is not valid JSON: %w", err)
	}
	jsonContentBytes, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	dir := filepath.Dir(yamlFilePath)
	jsonFilePath := filepath.Join(dir, "assignment.json")
	if err := os.WriteFile(jsonFilePath, jsonContentBytes, 0o644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	if err := os.Remove(yamlFilePath); err != nil {
		return fmt.Errorf("failed to delete original YAML file: %w", err)
	}
	return nil
}

// convertYamlKeysToJson converts YAML key-value pairs to JSON format by adding quotes around keys
func convertYamlKeysToJson(yamlContent string) string {
	lines := strings.Split(yamlContent, "\n")
	var jsonLines []string
	
	// Regex to match YAML key-value pairs (key: value)
	keyValueRegex := regexp.MustCompile(`^(\s*)([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*(.*)$`)
	
	for _, line := range lines {
		if keyValueRegex.MatchString(line) {
			// Replace key: value with "key": value
			matches := keyValueRegex.FindStringSubmatch(line)
			if len(matches) == 4 {
				indentation := matches[1]
				key := matches[2]
				value := matches[3]
				jsonLines = append(jsonLines, fmt.Sprintf(`%s"%s": %s`, indentation, key, value))
			} else {
				jsonLines = append(jsonLines, line)
			}
		} else {
			jsonLines = append(jsonLines, line)
		}
	}
	
	// Join with commas for JSON format
	return strings.Join(jsonLines, ",\n")
}
