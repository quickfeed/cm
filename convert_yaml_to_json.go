package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// convertYamlToJSON converts assignment.yml or assignment.yaml files to JSON format.
// It does this by simply wrapping the YAML content with { } and validating the result is valid JSON.
// If successful, it writes the content to assignment.json and deletes the original YAML file.
func convertYamlToJSON(args []string) {
	// Find all assignment.yml and assignment.yaml files in the current directory and subdirectories
	yamlFiles, err := findAssignmentYamlFiles(".")
	if err != nil {
		exitErr(err, "Error finding YAML files")
	}

	if len(yamlFiles) == 0 {
		fmt.Println("No assignment.yml or assignment.yaml files found.")
		return
	}

	converted := 0
	for _, yamlFile := range yamlFiles {
		fmt.Printf("Converting %s...\n", yamlFile)
		
		if err := convertSingleFile(yamlFile); err != nil {
			fmt.Printf("Error converting %s: %s\n", yamlFile, err)
			continue
		}
		
		converted++
		fmt.Printf("Successfully converted %s to assignment.json\n", yamlFile)
	}

	fmt.Printf("Converted %d file(s) successfully.\n", converted)
}

// findAssignmentYamlFiles finds all assignment.yml and assignment.yaml files in the given directory and its subdirectories
func findAssignmentYamlFiles(rootDir string) ([]string, error) {
	var yamlFiles []string
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		filename := info.Name()
		if filename == "assignment.yml" || filename == "assignment.yaml" {
			yamlFiles = append(yamlFiles, path)
		}
		
		return nil
	})
	
	return yamlFiles, err
}

// convertSingleFile converts a single YAML file to JSON
func convertSingleFile(yamlFilePath string) error {
	// Read the YAML file content
	content, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Trim whitespace from content
	contentStr := strings.TrimSpace(string(content))
	if contentStr == "" {
		return fmt.Errorf("YAML file is empty")
	}

	// Wrap content with { } to make it JSON-like
	jsonContent := "{" + contentStr + "}"

	// Validate that the result is valid JSON
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonContent), &jsonObj); err != nil {
		return fmt.Errorf("resulting content is not valid JSON: %w", err)
	}

	// Create the JSON file path in the same directory
	dir := filepath.Dir(yamlFilePath)
	jsonFilePath := filepath.Join(dir, "assignment.json")

	// Write the JSON content to the new file
	if err := os.WriteFile(jsonFilePath, []byte(jsonContent), 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	// Delete the original YAML file
	if err := os.Remove(yamlFilePath); err != nil {
		return fmt.Errorf("failed to delete original YAML file: %w", err)
	}

	return nil
}