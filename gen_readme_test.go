package main

import (
	"testing"
)

func TestStripMarkdownLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple markdown link",
			input:    "## [Docker Hub](https://hub.docker.com/) Setup",
			expected: "## Docker Hub Setup",
		},
		{
			name:     "multiple links",
			input:    "### Installing [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [GitHub Actions](https://github.com/features/actions)",
			expected: "### Installing Docker Desktop and GitHub Actions",
		},
		{
			name:     "no links",
			input:    "## Prerequisites",
			expected: "## Prerequisites",
		},
		{
			name:     "text with parentheses but no links",
			input:    "## Docker (Container Technology)",
			expected: "## Docker (Container Technology)",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripMarkdownLinks(tt.input)
			if result != tt.expected {
				t.Errorf("stripMarkdownLinks(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateToC(t *testing.T) {
	testMarkdown := `# Test Document

## [Docker Hub](https://hub.docker.com/) Setup

This section explains how to set up Docker Hub.

### Installing [Docker Desktop](https://www.docker.com/products/docker-desktop/)

Instructions for installing Docker Desktop.

## Prerequisites

Basic requirements.

### Working with [GitHub Actions](https://github.com/features/actions) CI/CD

How to use GitHub Actions.

## Summary

Final thoughts.`

	expectedHeadings := []string{
		"## Docker Hub Setup",
		"### Installing Docker Desktop",
		"## Prerequisites",
		"### Working with GitHub Actions CI/CD",
		"## Summary",
	}

	result := generateToC(testMarkdown)

	if len(result) != len(expectedHeadings) {
		t.Errorf("generateToC() returned %d headings, expected %d", len(result), len(expectedHeadings))
	}

	for i, expected := range expectedHeadings {
		if i < len(result) && result[i] != expected {
			t.Errorf("generateToC() heading %d = %q, want %q", i, result[i], expected)
		}
	}
}
