package main

import (
    "fmt"
    "os"

    "gommitizen/cmd"
    "gommitizen/internal"
    "gommitizen/internal/utils"
)

func main() {
    if len(os.Args) < 2 {
        cmd.HelpCommand()
        os.Exit(0)
    }

    // Parse the first argument (command), and the rest (flags/args)
    args := os.Args[1:]
    command := args[0]
    commandArgs := args[1:]

    // Handle help flags
    if command == "help" || command == "--help" || command == "-help" || command == "-h" {
        cmd.HelpCommand()
        os.Exit(0)
    }

    switch command {
    case "install":
        cmd.InstallCommand(commandArgs)
    case "reinstall":
        cmd.ReinstallCommand(commandArgs)
    case "uninstall":
        cmd.UninstallCommand(commandArgs)
    case "version":
        cmd.VersionCommand()
    case "commit":
        internal.CommitCommand(commandArgs)
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
        // Note: Lint options parsing expects args after "lint"
        opts, err := cmd.ParseLintOptions(commandArgs)
        if err != nil {
            fmt.Println("Failed to parse lint flags:", err)
            os.Exit(1)
        }

        if opts.All {
            err := internal.LintAllCommitMessage()
            if err != nil {
                fmt.Println(err.Error())
                os.Exit(1)
            }
            fmt.Println(utils.Color("All commit messages pass linting.", "green"))
        } else if opts.Current {
            err := internal.LintCurrentCommitMessage()
            if err != nil {
                fmt.Println(err.Error())
                os.Exit(1)
            }
            fmt.Println(utils.Color("Current commit message passes linting.", "green"))
        } else if opts.Message != "" {
            err := internal.LintSingleMessage(opts.Message)
            if err != nil {
                fmt.Println(err.Error())
                os.Exit(1)
            }
            fmt.Println(utils.Color("Provided message passes linting.", "green"))
        } else {
            fmt.Println("No lint target specified. Use --all, --current, or provide a message.")
            os.Exit(1)
        }
    default:
        fmt.Printf("Unknown command: %s\n\n", command)
        cmd.HelpCommand()
        os.Exit(1)
    }
}

