package main

import (
    "os"
    "fmt"
    "gommitizen/cmd"
    "gommitizen/internal"
)

func main() {
    if len(os.Args) < 2 {
        cmd.HelpCommand()
        os.Exit(0)
    }
    // Get the first argument.
    args := os.Args[1:]

    // If the user passes help flags or the help command, show help.
    if args[0] == "help" || args[0] == "--help" || args[0] == "-help" || args[0] == "-h" {
        cmd.HelpCommand()
        os.Exit(0)
    }

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
        // Assume commit message is provided as an argument.
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
        // If an unknown command is provided, show help.
        cmd.HelpCommand()
        os.Exit(1)
    }
}

