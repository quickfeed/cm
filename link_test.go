package main

import (
	"strings"
	"testing"
)

func TestLinkTemplateFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "heading text with markdown link",
			input:    "[Docker Hub](https://hub.docker.com/) Setup",
			expected: "docker-hub-setup",
		},
		{
			name:     "heading text without link",
			input:    "Prerequisites",
			expected: "prerequisites",
		},
		{
			name:     "heading text with multiple links",
			input:    "Installing [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [GitHub Actions](https://github.com/features/actions)",
			expected: "installing-docker-desktop-and-github-actions",
		},
		{
			name:     "heading text with special characters",
			input:    "Docker & Kubernetes: A Guide",
			expected: "docker-kubernetes-a-guide",
		},
		{
			name:     "header with multiple spaces and special chars",
			input:    "API & Database : Setup",
			expected: "api-database-setup",
		},
		{
			name:     "just text without header prefix",
			input:    "Docker Hub Setup",
			expected: "docker-hub-setup",
		},
	}

	// Get the link function from funcMap
	linkFunc := funcMap["link"].(func(string) string)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkFunc(tt.input)
			if result != tt.expected {
				t.Errorf("link(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCTemplateWithLinks(t *testing.T) {
	// Test the full template processing
	tmpl := parseTemplate("test", tocTemplate)

	labHeader := &LabHeader{
		LabHeader: "# Lab 4: Docker",
		ToC: []string{
			"## Table of Contents", // This should be skipped by slice .ToC 1
			"## [Docker Hub](https://hub.docker.com/) Setup",
			"### Installing [Docker Desktop](https://www.docker.com/products/docker-desktop/)",
			"## Prerequisites",
			"### Working with [GitHub Actions](https://github.com/features/actions) CI/CD",
		},
	}

	result := mustExecute(tmpl, labHeader)

	// Check that the result contains expected anchor links
	expectedLinks := []string{
		"#docker-hub-setup",
		"#installing-docker-desktop",
		"#prerequisites",
		"#working-with-github-actions-cicd",
	}

	for _, expectedLink := range expectedLinks {
		if !strings.Contains(result, expectedLink) {
			t.Errorf("Template output missing expected link: %s\nFull output:\n%s", expectedLink, result)
		}
	}

	// Check that markdown link syntax is removed from displayed text
	if strings.Contains(result, "[Docker Hub](https://hub.docker.com/)") {
		t.Errorf("Template output still contains markdown link syntax\nFull output:\n%s", result)
	}
}
