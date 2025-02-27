A lightweight, Go-based CLI tool for standardizing Git commit messages. It offers an interactive prompt that guides you through the process of composing a commit message that adheres to conventional commit formats. Gommitizen also includes commands for installing itself as a Git subcommand, checking version information, and more.

## Why Only the Standard Library?
Gommitizen is built exclusively with Go’s standard library. This design decision brings several benefits:

- Simplicity: Fewer dependencies means a cleaner, more maintainable codebase.
- Portability: The tool remains easily portable across different environments without worrying about external libraries.
- Reliability: Go’s standard library is well-tested and maintained, reducing the risk of introducing third-party bugs.
- Performance: Eliminating external dependencies minimizes build overhead and keeps execution fast.

## Features
- Interactive Commit Message Prompt: Guides you through entering the type, scope, subject, body, and footer of your commit message.
- Customizable Configuration: Uses an external JSON configuration (`configs/default.json`) to define form fields, options, and a commit message template.
- Git Integration: Executes Git commands to commit your changes using the generated commit message.
  - Self-Installation as Git Subcommand: Easily install Gommitizen as a Git subcommand (git-cz) with a single command.
  - Reinstall & Uninstall: Built-in commands to reinstall or uninstall the tool.

## Directory Structure

```
.
├── CHANGELOG.md
├── configs
│   └── default.json  # Default commit form configuration in JSON.
├── go.mod
├── internal
│   ├── changelog.go  # Changelog generation.
│   ├── cmd.go  # Contains install, reinstall, uninstall, and version commands.
│   ├── commit.go  # Handles commit message generation and execution.
│   ├── config.go  # Loads and renders configuration from configs/default.json.
│   ├── lint.go  # Commit message linter.
│   └── version.go  # Version management.
├── LICENSE
├── main.go  # Main entry point that dispatches commands.
└── README.md
```

## Installation

### Clone the Repository:

```
git clone https://github.com/tiendu/gommitizen.git
cd gommitizen
```

### Build the Binary:

- Use the Go compiler to build the tool: `go build -o gommitizen`.

- This produces a portable binary named gommitizen.

## Usage

Gommitizen supports several subcommands. The main usage pattern is: `./gommitizen <command> [options]`

### Commands
- `install`: Installs the tool as a Git subcommand (git-cz).
- `reinstall`: Uninstalls and then reinstalls the tool.
- `uninstall`: Uninstalls the tool from Git's exec path.
- `version`: Prints version information.
- `commit`: Runs the interactive commit prompt, which loads configuration from `configs/default.json`, collects user input, renders the commit message, and then executes git commit.
- `lint`: Runs a linter to validate the commit message.
- `changelog`: Generates a changelog depending on the type of commit.
- `bump`: Updates the VERSION with the current version.
- `COMMITIZEN_INSTALL_PATH`: Custom path for the installation.

### Configuration

Gommitizen uses a JSON configuration file located at `configs/default.json`. This file defines:
- Form Fields: Such as type, scope, subject, body, and footer
- Options for "select" Fields: For example, the available types of changes like feat, fix, docs, etc.
- Commit Message Template: A Go template that assembles the final commit message based on your input.

You can modify `configs/default.json` to customize the prompts and commit message format.

