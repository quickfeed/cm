// Package main provides the cm (course management) command-line tool for QuickFeed.
//
// The cm tool helps instructors manage course repositories, assignments, and
// student repositories for courses using the QuickFeed system. It provides
// commands for initializing course environments, managing repositories,
// generating documentation, and maintaining course materials.
//
// # Installation
//
// To install the cm tool:
//
//	go install github.com/quickfeed/cm@latest
//
// # Usage
//
// The general usage pattern is:
//
//	cm <command> [options]
//
// # Common Workflows
//
// Setting up a new course:
//
//	cm init-env -year 2025 -course dat520 -name "Distributed Systems"
//	cm init-repos
//
// Managing student repositories:
//
//	cm list-repos
//	cm clone-repo -repo username-labs
//	cm clone-all-repos
//
// Generating course materials:
//
//	cm gen-readme
//	cm gen-tests-json -labs lab1
//	cm update-doc-tags -repo assignments
//
// # Commands
//
// Use 'cm help' to see all available commands, or 'cm help <command>' for
// detailed help on a specific command.
//
// # Environment Variables
//
// The tool uses environment variables stored in a .env file in the repository
// root. These are typically set up using the 'cm init-env' command:
//
//	YEAR - Course year (e.g., 2025)
//	COURSE - Course code (e.g., dat520)
//	NAME - Course name (e.g., "Distributed Systems")
//	COURSE_ORG - GitHub organization (defaults to COURSE-YEAR)
//	BOT_USER - Help bot username (defaults to "helpbot")
//	DISCORD_JOIN_URL - Discord server join URL (optional)
//
// # Repository Structure
//
// The tool expects a specific repository structure:
//
//	<course-repo>/           # Main course repository
//	├── .env                 # Environment variables
//	├── assignments/         # Assignment templates
//	├── info/                # Course information
//	├── tests/               # Test files
//	└── <year>/              # Year-specific repositories
//	    ├── assignments/     # Published assignments
//	    ├── info/            # Published course info
//	    ├── tests/           # Published tests
//	    └── student-repos/   # Cloned student repositories
//
// For more information, see the QuickFeed documentation.
package main
