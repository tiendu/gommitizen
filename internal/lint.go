package internal

import (
    "fmt"
    "strings"
    "math"
    "os"
    "io"
    "os/exec"
    "regexp"
    "sync"
)

// calculateEntropy computes the Shannon entropy of a string.
func calculateEntropy(s string) float64 {
    freq := make(map[rune]float64)
    for _, r := range s {
        freq[r]++
    }
    var entropy float64
    l := float64(len(s))
    for _, count := range freq {
        p := count / l
        entropy -= p * math.Log2(p)
    }
    return entropy
}

// LintSensitiveFiles checks staged files concurrently for sensitive information,
// ignoring files in the "internal/" directory.
func LintSensitiveFiles() error {
    cmd := exec.Command("git", "diff", "--cached", "--name-only")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to get staged files: %v", err)
    }

    files := strings.Split(strings.TrimSpace(string(output)), "\n")
    if len(files) == 0 || files[0] == "" {
        return nil
    }

    // Use word boundaries in regex patterns.
    patterns := []string{
        `(?i)\bpassword\b`,
        `(?i)\bsecret\b`,
        `(?i)\bapikey\b`,
        `(?i)\btoken\b`,
        `(?i)\bcredential\b`,
    }

    var regexes []*regexp.Regexp
    for _, pat := range patterns {
        re, err := regexp.Compile(pat)
        if err != nil {
            return fmt.Errorf("failed to compile regex %s: %v", pat, err)
        }
        regexes = append(regexes, re)
    }

    var wg sync.WaitGroup
    errCh := make(chan error, len(files))
    for _, file := range files {
        // Skip files in the "internal/" directory.
        if strings.HasPrefix(file, "internal/") {
            continue
        }

        // Open file.
        f, err := os.Open(file)
        if err != nil {
            fmt.Printf("Warning: unable to open file %s: %v\n", Color(file, "yellow"), err)
            continue
        }

        wg.Add(1)
        // Pass both the filename and the file handle to the goroutine.
        go func(file string, f *os.File) {
            defer wg.Done()
            defer f.Close()
            data, err := io.ReadAll(f)
            if err != nil {
                // Skip files that cannot be read.
                return
            }

            for _, re := range regexes {
                if re.Match(data) {
                    // Extract the matched substring.
                    match := re.Find(data)
                    entropy := calculateEntropy(string(match))
                    // Only flag if entropy exceeds a threshold (e.g., 3.5).
                    if entropy > 3.5 {
                        errCh <- fmt.Errorf("sensitive information detected in file: %s (pattern: %s, entropy: %.2f)", Color(file, "yellow"), re.String(), entropy)
                        return
                    }
                }
            }
        }(file, f)
    }

    wg.Wait()
    close(errCh)

    var combinedErrs []string
    for err := range errCh {
        combinedErrs = append(combinedErrs, err.Error())
    }
    if len(combinedErrs) > 0 {
        return fmt.Errorf("linting failed:\n%s", strings.Join(combinedErrs, "\n"))
    }

    return nil
}

// LintCommitMessage checks if the commit message adheres to predefined rules.
func LintCommitMessage(message string) error {
    // Example: Ensure the subject (first line) is non-empty and under 100 characters.
    lines := strings.SplitN(message, "\n", 2)
    subject := strings.TrimSpace(lines[0])
    if len(subject) == 0 {
        return fmt.Errorf("%s", "commit subject cannot be empty")
    }
    if len(subject) > 100 {
        return fmt.Errorf("%s", "commit subject is too long (max 100 characters)")
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

// LintCurrentCommitMessage lints the current commit messages.
func LintCurrentCommitMessage() error {
    // Get the latest commit message from HEAD.
    cmd := exec.Command("git", "log", "-1", "--pretty=%B")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to get commit message: %v", err)
    }
    message := strings.TrimSpace(string(output))
    // Lint the commit message using LintCommitMessage.
    return LintCommitMessage(message)
}

// LintAllCommitMessage lints all commit messages.
func LintAllCommitMessage() error {
    // Get all commit hashes from the repository.
    cmd := exec.Command("git", "log", "--pretty=%H")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to get commit hashes: %v", err)
    }
    hashes := strings.Split(strings.TrimSpace(string(output)), "\n")
    if len(hashes) == 0 {
        return nil
    }

    var combinedErrors []string
    // Iterate over each commit hash.
    for _, hash := range hashes {
        // Get the commit message for this commit.
        cmdMsg := exec.Command("git", "log", "-1", "--pretty=%B", hash)
        msgOut, err := cmdMsg.Output()
        if err != nil {
            combinedErrors = append(combinedErrors, fmt.Sprintf("failed to get commit message for %s: %v", Color(hash, "red"), err))
            continue
        }
        message := strings.TrimSpace(string(msgOut))
        if err := LintCommitMessage(message); err != nil {
            combinedErrors = append(combinedErrors, fmt.Sprintf("commit %s: %v", Color(hash, "red"), err))
        }
    }
    if len(combinedErrors) > 0 {
        return fmt.Errorf("linting errors found:\n%s", strings.Join(combinedErrors, "\n"))
    }
    return nil
}

