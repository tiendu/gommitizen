package main

import (
    "os"
    "fmt"

    "gommitizen/cmd"
    "gommitizen/internal"
    "gommitizen/internal/utils"
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
        // Parse lint-specific flags using the helper.
        opts, err := cmd.ParseLintOptions(os.Args[2:])
        if err != nil {
            fmt.Println("Failed to parse lint flags:", err)
            os.Exit(1)
        }

        if opts.All {
            err = internal.LintAllCommitMessage()
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
        }
    default:
        // If an unknown command is provided, show help.
        cmd.HelpCommand()
        os.Exit(1)
    }
}

