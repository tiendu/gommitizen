package main

import (
    "fmt"
    "os"
    "gommitizen/internal"
    "gommitizen/cmd"
)

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }
    args := os.Args[1:]

    switch args[0] {
    case "install":
        cmd.InstallCommand()
    case "reinstall":
        cmd.ReinstallCommand()
    case "uninstall":
        cmd.UninstallCommand()
    case "version":
        cmd.VersionCommand()
    case "commit":
        internal.CommitCommand(args)
    case "changelog":
        if err := internal.GenerateChangelog(); err != nil {
            fmt.Printf("Changelog generation failed: %v\n", err)
        }
    case "bump":
        if newVersion, err := cmd.BumpVersion(); err != nil {
            fmt.Printf("Version bump failed: %v\n", err)
        } else {
            fmt.Printf("New version: %s\n", newVersion)
        }
    case "lint":
        // For demonstration, assume commit message is provided as an argument.
        if len(args) < 1 {
            fmt.Println("Usage: gommitizen lint <commit-message>")
            os.Exit(1)
        }
        message := args[0]
        if err := internal.LintCommitMessage(message); err != nil {
            fmt.Printf("Lint failed: %v\n", err)
            os.Exit(1)
        } else {
            fmt.Println("Commit message passes linting.")
        }
    default:
        fmt.Printf("Unknown command: %s\n", args[1])
        usage()
        os.Exit(1)
    }
}

func usage() {
    fmt.Println("Usage: gommitizen <command> [options]")
    fmt.Println("Commands:")
    fmt.Println("  install      Install this tool as a git subcommand (git-cz)")
    fmt.Println("  version      Print version information")
    fmt.Println("  commit       Create a commit using the configured commitizen flow")
    fmt.Println("  changelog    Generate a CHANGELOG.md from commit logs")
    fmt.Println("  bump         Bump the version automatically")
    fmt.Println("  lint         Lint a commit message")
}

