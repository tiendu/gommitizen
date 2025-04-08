A lightweight, Go-based configurable CLI tool for standardizing Git commit messages. It offers an interactive prompt that guides you through the process of composing a commit message that adheres to conventional commit formats. Gommitizen also includes commands for installing itself as a Git subcommand, checking version information, and more.

## Directory Structure

```
.
├── CHANGELOG               # Project changelog
├── cmd
│   ├── options.go          # Install, reinstall, uninstall, and version commands
│   └── version.go          # Version management
├── configs
│   └── default.json        # Default commit form configuration (JSON)
├── go.mod                  # Go module definition
├── gommitizen              # Built binary (after build)
├── internal
│   ├── changelog.go        # Changelog generation
│   ├── commit.go           # Commit message generation and execution
│   ├── config.go           # Load and render config from configs/default.json
│   ├── lint.go             # Commit message linter
│   └── utils               # Terminal UI and utilities
│       ├── term_darwin.go  # Terminal handling for macOS
│       ├── term_linux.go   # Terminal handling for Linux
│       ├── text_mods.go    # Text formatting utilities (colors, underline, highlight)
│       └── tui.go          # Terminal UI (TUI) core logic
├── LICENSE                 # License file (MIT)
├── main.go                 # Application entry point
├── README.md               # Project documentation
└── VERSION                 # Current app version
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

## Usage

After installation, you'll have `git-cz` command available globally.

### Create a Commit

```bash
git-cz commit
```

Launches an interactive prompt to compose your commit message.

### Lint Commit Messages

```bash
# Lint all commit messages in the repository
git-cz lint --all
git-cz lint -a

# Lint only the latest commit message
git-cz lint --current
git-cz lint -c

# Optionally, lint a specific commit message string directly
git-cz lint "your commit message here"
```

### Install / Reinstall / Uninstall

```bash
git-cz install      # Install gommitizen
git-cz reinstall    # Reinstall gommitizen
git-cz uninstall    # Uninstall gommitizen
```

## Configuration

Upon installation, gommitizen copies the config to:

```bash
~/.local/bin/gommitizen/configs/default.json
```

You can edit this file to customize:

- Commit types
- Prompts
- Commit message templates

If no config file is found, gommitizen will use its built-in default config.

