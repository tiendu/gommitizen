package internal

import (
    "fmt"
    "strings"
)

// LintCommitMessage checks if the commit message adheres to predefined rules.
func LintCommitMessage(message string) error {
    // Example: Ensure the subject (first line) is non-empty and under 100 characters.
    lines := strings.SplitN(message, "\n", 2)
    subject := strings.TrimSpace(lines[0])
    if len(subject) == 0 {
        return fmt.Errorf("commit subject cannot be empty")
    }
    if len(subject) > 100 {
        return fmt.Errorf("commit subject is too long (max 100 characters)")
    }

    // Example: Optionally enforce that the commit message starts with a valid type.
    validTypes := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "revert", "WIP"}
    valid := false
    for _, t := range validTypes {
        if strings.HasPrefix(subject, t) {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("commit subject must start with one of the following types: %v", validTypes)
    }
    return nil
}


