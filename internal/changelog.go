package internal

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
)

// commitEntry represents a parsed commit.
type commitEntry struct {
    hash    string
    date    string
    ctype   string
    subject string
}

// GenerateChangelog runs "git log" to extract commit messages (including commit date),
// processes them concurrently, groups them by commit type, and writes the results to CHANGELOG.md.
func GenerateChangelog() error {
    // Run git log with a custom format: hash | date | subject.
    // Using a delimiter (|) makes parsing easier.
    cmd := exec.Command("git", "log", "--pretty=format:%h|%ad|%s", "--date=iso")
    out, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to run git log: %v", err)
    }

    // Split the log output into individual commit lines.
    lines := strings.Split(string(out), "\n")
    if len(lines) == 0 {
        return fmt.Errorf("no commits found")
    }

    // Channel to collect parsed commit entries.
    commitCh := make(chan commitEntry, len(lines))
    var wg sync.WaitGroup

    // Process each commit line concurrently.
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        wg.Add(1)
        go func(l string) {
            defer wg.Done()
            // Expect the format: "<hash>|<date>|<subject>"
            parts := strings.SplitN(l, "|", 3)
            if len(parts) < 3 {
                // Skip malformed lines.
                return
            }
            hash := parts[0]
            date := strings.TrimSpace(parts[1])
            subjectLine := parts[2]

            // Extract commit type from subject (if a colon exists).
            ctype := "unknown"
            if colonIdx := strings.Index(subjectLine, ":"); colonIdx != -1 {
                ctype = strings.TrimSpace(subjectLine[:colonIdx])
            }

            commitCh <- commitEntry{
                hash:    hash,
                date:    date,
                ctype:   ctype,
                subject: subjectLine,
            }
        }(line)
    }

    // Wait for all goroutines to finish processing.
    wg.Wait()
    close(commitCh)

    // Group commits by type.
    groups := make(map[string][]commitEntry)
    for entry := range commitCh {
        groups[entry.ctype] = append(groups[entry.ctype], entry)
    }

    // Build the changelog content.
    var buf bytes.Buffer
    buf.WriteString("# Changelog\n\n")
    for typ, entries := range groups {
        buf.WriteString(fmt.Sprintf("## %s\n\n", typ))
        for _, e := range entries {
            // Include hash, date, and subject for each commit.
            buf.WriteString(fmt.Sprintf("- [%s] %s %s\n", e.hash, e.date, e.subject))
        }
        buf.WriteString("\n")
    }

    // Write the content to CHANGELOG.md.
    changelogPath := filepath.Join(".", "CHANGELOG.md")
    if err := ioutil.WriteFile(changelogPath, buf.Bytes(), 0644); err != nil {
        return fmt.Errorf("failed to write CHANGELOG.md: %v", err)
    }
    fmt.Println("Changelog generated in", changelogPath)
    return nil
}

