package main

import (
    "os"
    "gommitizen/cmd"
)

func main() {
    if len(os.Args) < 2 {
        cmd.HelpCommand()
        os.Exit(0)
    }
    // Get the first argument.
    arg := os.Args[1]

    // If the user passes help flags or the help command, show help.
    if arg == "help" || arg == "--help" || arg == "-help" || arg == "-h" {
        cmd.HelpCommand()
        os.Exit(0)
    }

    switch arg {
    case "install":
        cmd.InstallCommand()
    case "reinstall":
        cmd.ReinstallCommand()
    case "uninstall":
        cmd.UninstallCommand()
    case "version":
        cmd.VersionCommand()
    // Add other commands like commit, changelog, bump, lint, etc.
    default:
        // If an unknown command is provided, show help.
        cmd.HelpCommand()
        os.Exit(1)
    }
}

