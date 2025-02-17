package internal

import (
    "fmt"
    "strings"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "sync"
)

// LintSensitiveFilesConcurrent checks staged files concurrently.
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

    patterns := []string{
        "(?i)password",
        "(?i)secret",
        "(?i)apikey",
        "(?i)token",
        "(?i)credential",
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
            // If a file can't be opened, log and continue.
            fmt.Printf("Warning: unable to open file %s: %v\n", file, err)
            continue
        }

        defer f.Close()

        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            data, err := ioutil.ReadFile(f)
            if err != nil {
                // Skip files that cannot be read.
                return
            }

            for _, re := range regexes {
                if re.Match(data) {
                    errCh <- fmt.Errorf("sensitive information detected in file: %s (pattern: %s)", f, re.String())
                    return
                }
            }
        }(file)
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

