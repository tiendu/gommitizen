package main

import (
    "fmt"
    "os"
    "gommitizen/internal"
)

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

    // The first argument is the command.
    cmd := os.Args[1]
    // The rest are arguments for that command.
    args := os.Args[2:]

    switch cmd {
    case "install":
        internal.InstallCommand()
    case "reinstall":
        internal.ReinstallCommand()
    case "uninstall":
        internal.UninstallCommand()
    case "version":
        internal.VersionCommand()
    case "commit":
        internal.CommitCommand(args)
    default:
        fmt.Printf("Unknown command: %s\n", cmd)
        usage()
        os.Exit(1)
    }
}

func usage() {
    fmt.Println("Usage: gommitizen <command> [options]")
    fmt.Println("Commands:")
    fmt.Println("  install    Install this tool as a git subcommand (git-cz)")
    fmt.Println("  version    Print version information")
    fmt.Println("  commit     Create a commit using the configured commitizen flow")
}

