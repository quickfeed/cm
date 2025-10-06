# Course Management Tool (cm)

The `cm` (course management) tool is a command-line utility designed for QuickFeed instructors to manage course repositories, assignments, and student repositories. It provides commands for initializing course environments, managing repositories, generating documentation, and maintaining course materials.

## Table of Contents

- [Installation](#installation)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Initial Setup](#initial-setup)
- [Common Workflows](#common-workflows)
- [Environment Variables](#environment-variables)
- [Repository Structure](#repository-structure)
- [Available Commands](#available-commands)
- [Using Just Commands](#using-just-commands)
- [Troubleshooting](#troubleshooting)

## Installation

Before installing the `cm` tool, you need to configure the `GOPRIVATE` environment variable to access the private repository:

```bash
# Set GOPRIVATE to allow access to the cm repository
export GOPRIVATE=github.com/quickfeed/cm

# Add this to your shell configuration file (~/.bashrc, ~/.zshrc, etc.)
echo 'export GOPRIVATE=github.com/quickfeed/cm' >> ~/.bashrc  # or ~/.zshrc
```

Then install the `cm` tool:

```bash
go install github.com/quickfeed/cm@latest
```

Alternatively, if you're working within a course repository, you can use the `just` command to set up the tool:

```bash
just setup
```

This creates a `tools.mod` file to manage the `cm` tool version separately from the main `go.mod` file.

## Prerequisites

The following tools are required for working with course repositories:

- **Go** (1.21 or later) - for running the `cm` tool
- **GitHub CLI (`gh`)** - for GitHub authentication and API access
- **Just** - for running automation recipes
- **Tree** - for viewing directory structures
- **golangci-lint** - for linting Go code (optional, for development)
- **Protocol Buffers (`protobuf`)** - for certain course assignments (optional)

You can install these tools on macOS using Homebrew:

```bash
brew install protobuf gh just tree golangci-lint
```

Or use the provided `just` command:

```bash
just tools
```

## Quick Start

For course instructors setting up a new course:

```bash
# 1. Install prerequisites
just tools

# 2. Set up the cm tool
just setup

# 3. Create environment configuration
just env

# 4. Initialize course repositories
just init
```

## Initial Setup

### Step 1: Set Up the cm Tool

If you're starting with a new course repository:

```bash
# Set up cm tool with tools.mod file
just setup
```

This command creates a `tools.mod` file for managing the `cm` tool separately.

### Step 2: Create Environment Configuration

Initialize the `.env` file with your course information:

```bash
# Using the cm tool directly
cm init-env -year 2025 -course dat520 -name "Distributed Systems"

# Or using just
just env  # Edit the Justfile first to set your course details
```

The `.env` file will contain required environment variables like `YEAR`, `COURSE`, `NAME`, and optionally `DISCORD_JOIN_URL` and `BOT_USER`.

**Important:** Commit the `.env` file to your repository and update it as needed for your course.

### Step 3: Initialize Course Repositories

Create the necessary course repositories on GitHub:

```bash
# Using the cm tool
cm init-repos

# Or using just
just init
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

### Syncing Changes Using Just

```bash
# Update README files for all labs
just readme

# Sync changes to the info repository
just info

# Sync changes to assignments (specify lab folders)
just assignments lab1 lab2

# Sync changes to tests (specify lab folders)
just tests lab1 lab2
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
| `GOPRIVATE` | Go private module configuration | Auto-set | `github.com/quickfeed/cm` |

The `GOPRIVATE` variable is automatically added to your `.env` file when you run `cm init-env`.

## Repository Structure

The `cm` tool expects a specific repository structure:

```
<course-repo>/           # Main course repository (internal)
├── .env                 # Environment variables
├── Justfile             # Automation recipes
├── tools.mod            # Go tools configuration
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

## Using Just Commands

The `Justfile` provides convenient recipes for common tasks. Use `just --list` to see all available commands.

### Common Just Recipes

| Recipe | Description |
|--------|-------------|
| `just setup` | Set up the cm tool using tools.mod |
| `just update-cm` | Update the cm tool to the latest version |
| `just tools` | Install required tools (macOS/Homebrew) |
| `just env` | Create .env file (edit Justfile first) |
| `just init` | Initialize course repositories |
| `just readme` | Update README.md files for all labs |
| `just clone` | Clone your personal labs repository |
| `just info` | Sync changes to info repository |
| `just assignments lab1` | Sync changes to assignments for lab1 |
| `just tests lab1` | Sync changes to tests for lab1 |

### Workflow Example

```bash
# Initial setup
just setup
just env
just init

# Working with assignments
just readme
just assignments lab1 lab2
just tests lab1 lab2

# Sync course information
just info
```

## Troubleshooting

### GOPRIVATE Not Set

If you get errors about accessing the private repository:

```bash
# Ensure GOPRIVATE is set
export GOPRIVATE=github.com/quickfeed/cm

# Add to your shell configuration
echo 'export GOPRIVATE=github.com/quickfeed/cm' >> ~/.zshrc  # or ~/.bashrc
```

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
# Update using just
just update-cm

# Or update manually
go get -tool -modfile=tools.mod github.com/quickfeed/cm@latest
```

---

For more information, see the [QuickFeed documentation](https://github.com/quickfeed/quickfeed).
