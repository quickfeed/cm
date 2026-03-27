# Course Management Tool (cm)

The `cm` (course management) tool is a command-line utility designed for QuickFeed instructors to manage course repositories, assignments, and student repositories. It provides commands for initializing course environments, managing repositories, generating documentation, and maintaining course materials.

## Table of Contents

- [Installation](#installation)
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Common Workflows](#common-workflows)
- [Environment Variables](#environment-variables)
- [Repository Structure](#repository-structure)
- [Available Commands](#available-commands)
- [Troubleshooting](#troubleshooting)

## Installation

To install the `cm` tool:

```bash
go install github.com/quickfeed/cm@latest
```

## Prerequisites

The following tools are required for working with the `cm` tool:

- **Go** (1.24 or later) - for running the `cm` tool
- **GitHub CLI (`gh`)** - for GitHub authentication and API access

## Initial Setup

### Step 1: Create Environment Configuration

Initialize the `.env` file with your course information:

```bash
cm init-env -year 2025 -course dat520 -name "Distributed Systems"
```

The `.env` file will contain required environment variables like `YEAR`, `COURSE`, `NAME`, and optionally `DISCORD_JOIN_URL` and `BOT_USER`.

**Important:** Commit the `.env` file to your repository and update it as needed for your course.

### Step 2: Initialize Course Repositories

Create the necessary course repositories on GitHub:

```bash
cm init-repos
```

This will create three repositories in your course organization:
- `assignments` - Published assignments for students
- `info` - Course information and documentation
- `tests` - Test files for assignments

## Common Workflows

### Setting Up a New Course

```bash
cm init-env -year 2025 -course dat520 -name "Distributed Systems"
cm init-repos
```

### Managing Student Repositories

```bash
# List all student repositories
cm list-repos

# Clone a specific student repository
cm clone-repo -repo username-labs

# Clone all student repositories
cm clone-all-repos
```

### Generating Course Materials

```bash
# Generate README.md files for all labs
cm gen-readme

# Generate tests.json for specific labs
cm gen-tests-json -labs lab1

# Update documentation tags in repositories
cm update-doc-tags -repo assignments
```

## Environment Variables

The `cm` tool uses environment variables stored in a `.env` file in the repository root. These are typically set up using the `cm init-env` command:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `YEAR` | Course year (e.g., 2025) | Yes | - |
| `COURSE` | Course code (e.g., dat520) | Yes | - |
| `NAME` | Course name (e.g., "Distributed Systems") | Yes | - |
| `COURSE_ORG` | GitHub organization name | No | `{COURSE}-{YEAR}` |
| `BOT_USER` | Help bot username | No | `helpbot` |
| `DISCORD_JOIN_URL` | Discord server join URL | No | - |

## Repository Structure

The `cm` tool expects a specific repository structure:

```
<course-repo>/           # Main course repository (internal)
├── .env                 # Environment variables
├── assignments/         # Assignment templates
├── info/                # Course information templates
├── tests/               # Test file templates
└── <year>/              # Year-specific published repositories
    ├── assignments/     # Published assignments (public/student-facing)
    ├── info/            # Published course info (public/student-facing)
    ├── tests/           # Published tests (public/student-facing)
    └── student-repos/   # Cloned student repositories
```

**Note:** The main course repository contains templates and internal materials, while the `<year>/` subdirectories contain published, student-facing materials.

## Available Commands

Use `cm help` to see all available commands, or `cm help <command>` for detailed help on a specific command.

| Command | Description |
|---------|-------------|
| `init-env` | Initialize environment variables for the course |
| `init-repos` | Initialize course repositories (assignments, info, tests) |
| `list-repos` | List all student and group repositories |
| `clone-repo` | Clone a specified repository or current user's repository |
| `clone-all-repos` | Clone all student and group repositories |
| `gen-readme` | Generate README.md files from readme_tmpl.md templates |
| `gen-tests-json` | Generate tests.json files for lab assignments |
| `update-doc-tags` | Update documentation tags in markdown files |
| `remove-solution-tags` | Remove solution build tags from Go files |
| `rename-tests` | Rename legacy `*_ag_test.go` files to `*_qf_test.go` |
| `add-lint-checkers` | Add linter test files to specified lab folder |
| `add-main-tests` | Add main test files to directories with existing test files |
| `convert-yaml-to-json` | Convert assignment.yml/yaml files to assignment.json |

### Example Usage

```bash
# Initialize environment
cm init-env -year 2025 -course dat520 -name "Distributed Systems"

# List repositories in your course organization
cm list-repos

# Generate README files
cm gen-readme

# Generate tests.json for multiple labs
cm gen-tests-json -labs "lab1 lab2 lab3"

# Update documentation tags
cm update-doc-tags -repo assignments
```

## Troubleshooting

### GitHub Authentication

Ensure you're authenticated with the GitHub CLI:

```bash
gh auth login
gh auth status
```

### Missing .env File

If commands fail with missing environment variables:

```bash
# Create the .env file
cm init-env -year 2025 -course dat520 -name "Distributed Systems"

# Verify the file was created
cat .env
```

### Permission Issues

If you get permission errors when initializing repositories, ensure:
- You have admin access to the GitHub organization
- The organization name matches your `COURSE_ORG` environment variable
- You're authenticated with GitHub CLI (`gh auth status`)

### Tool Version Issues

To update to the latest version of the cm tool:

```bash
go install github.com/quickfeed/cm@latest
```

---

For more information, see the [QuickFeed documentation](https://github.com/quickfeed/quickfeed).
