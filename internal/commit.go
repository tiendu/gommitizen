package internal

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

// isGitRepo checks if the current directory is a Git repository.
func isGitRepo() bool {
    cwd, err := os.Getwd()
    if err != nil {
        return false
    }
    gitPath := filepath.Join(cwd, ".git")
    info, err := os.Stat(gitPath)
    if err != nil {
        return false
    }
    if info.IsDir() {
        return true
    }
    // If .git is a file, check if it contains a reference to the actual git directory.
    data, err := os.ReadFile(gitPath)
    if err != nil {
        return false
    }
    return strings.HasPrefix(string(data), "gitdir:")
}

// CommitCommand loads configuration, collects user input, renders the commit message, and executes the git commit command.
func CommitCommand(args []string) {
    commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
    allFlag := commitFlags.Bool("all", false, "Automatically stage modified/deleted files")
    commitFlags.Parse(args)

    if !isGitRepo() {
        fmt.Println("Current directory is not a git repository.")
        return
    }

    // Load configuration from external file.
    config, err := LoadConfig()
    if err != nil {
        log.Printf("No external config loaded: %v; using built-in default\n", err)
        config = LoadDefaultConfig()
    }

    // Collect user input based on the configuration.
    answers := CollectUserInput(config)

    // Render the commit message template using the collected answers.
    message, err := RenderTemplate(config, answers)
    if err != nil {
        log.Printf("Error rendering commit message: %v\n", err)
        return
    }

    // Within CommitCommand after rendering the message:
    if err := LintCommitMessage(message); err != nil {
        log.Printf("Commit message linting failed: %v\n", err)
        return
    }

    if err := LintSensitiveFiles(); err != nil {
        log.Printf("%v\n", err)
        return
    }

    // Execute git commit with the assembled message.
    output, err := commitMessage(message, *allFlag)
    if err != nil {
        log.Printf("Git commit failed: %v\n", err)
        log.Printf("Commit message was:\n%s\n", message)
    }

    fmt.Print(output)
}

// commitMessage executes the "git commit" command with the given message.
func commitMessage(message string, all bool) (string, error) {
    args := []string{"commit"}
    if all {
        args = append(args, "-a")
    }

    args = append(args, "-m", message)
    cmd := exec.Command("git", args...)
    out, err := cmd.CombinedOutput()
    return string(out), err
}

