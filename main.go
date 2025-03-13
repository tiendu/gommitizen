package main

import (
    "os"
    "fmt"
    "flag"

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
    lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
    // Define both long and short forms.
    lintAllLong := lintFlags.Bool("all", false, "Lint all commit messages")
    lintAllShort := lintFlags.Bool("a", false, "Lint all commit messages (short)")
    lintCurrentLong := lintFlags.Bool("current", false, "Lint the current commit message")
    lintCurrentShort := lintFlags.Bool("c", false, "Lint the current commit message (short)")
    lintFlags.Parse(os.Args[2:])

    // Combine flag values: if either the long or short flag is true, treat it as active.
    lintAll := *lintAllLong || *lintAllShort
    lintCurrent := *lintCurrentLong || *lintCurrentShort

    if lintAll {
        err := internal.LintAllCommitMessage()
        if err != nil {
            fmt.Printf("Lint all commit messages failed:\n%v\n", err)
            os.Exit(1)
        }
        fmt.Println("All commit messages pass linting.")
    } else if lintCurrent {
        err := internal.LintCurrentCommitMessage()
        if err != nil {
            fmt.Printf("Current commit message linting failed:\n%v\n", err)
            os.Exit(1)
        }
        fmt.Println("Current commit message passes linting.")
    } else {
        if lintFlags.NArg() < 1 {
            fmt.Println("Usage: gommitizen lint --all (or -a) OR gommitizen lint --current (or -c)")
            os.Exit(1)
        }
        // Optionally, lint a provided commit message directly.
        message := lintFlags.Arg(0)
        if err := internal.LintCommitMessage(message); err != nil {
            fmt.Printf("Lint failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Commit message passes linting.")
    }
    default:
        // If an unknown command is provided, show help.
        cmd.HelpCommand()
        os.Exit(1)
    }
}

