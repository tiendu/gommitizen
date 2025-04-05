A lightweight, Go-based configurable CLI tool for standardizing Git commit messages. It offers an interactive prompt that guides you through the process of composing a commit message that adheres to conventional commit formats. Gommitizen also includes commands for installing itself as a Git subcommand, checking version information, and more.

## Directory Structure

```
.
├── CHANGELOG.md
├── cmd
│   ├── options.go  # Contains install, reinstall, uninstall, and version commands.
│   └── version.go  # Version management.
├── configs
│   └── default.json  # Default commit form configuration in JSON.
├── go.mod
├── internal
│   ├── changelog.go  # Changelog generation.
│   ├── commit.go # Handles commit message generation and execution.
│   ├── config.go  # Loads and renders configuration from configs/default.json.
│   ├── lint.go  # Commit message linter.
├── LICENSE
├── main.go
└── README.md
```

## Installation

### Clone the Repository:

```
git clone https://github.com/tiendu/gommitizen.git
cd gommitizen
```

### Default install (recommended)

This will install `gommitizen` into `~/.local/bin/gommitizen` and copy the config.

```bash
go build -o gommitizen
./gommitizen install
```

### Custom install path

You can specify a custom installation path:

```bash
./gommitizen install --path /your/custom/path
```

**NOTE**:  Make sure to add your install path to your `PATH` environment variable!

## Commands

### Install

Install the tool to `~/.local/bin/gommitizen` and copy default configs.

### Reinstall

Uninstall and reinstall the tool.

### Uninstall

### Commit

Create a commit interactively with `git-cz` command using the configured commit flow.

### Lint

Lint commit messages.

```bash
# Lint all commit messages in the repository
git-cz lint --all
git-cz lint -a

# Lint only the current (latest) commit message
git-cz lint --current
git-cz lint -c

# Optionally, lint a specific commit message directly
git-cz lint "your commit message here"
```

## Configuration

When installed, gommitizen copies its config to: `~/.local/bin/gommitizen/configs/default.json`

You can edit this file to customize:

- Commit types
- Prompts
- Template for commit messages

If no config file is found, `gommitizen` will use a built-in default config.



