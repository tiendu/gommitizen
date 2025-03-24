package internal

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
)

// commitEntry represents a parsed commit.
type commitEntry struct {
    hash    string
    date    string
    author  string
    ctype   string
    subject string
}

// GenerateChangelog runs "git log" to extract commit messages (including commit date and author),
// processes them concurrently, groups them by commit type, and writes the results to CHANGELOG.md.
func GenerateChangelog() error {
    // Run git log with a custom format: hash | date | author | subject.
    // Using a delimiter (|) makes parsing easier.
    // --date=iso will output the commit date in ISO 8601 format (which includes the timezone offset).
    cmd := exec.Command("git", "log", "--pretty=format:%h|%ad|%an|%s", "--date=iso")
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
            // Expect the format: "<hash>|<date>|<author>|<subject>"
            parts := strings.SplitN(l, "|", 4)
            if len(parts) < 4 {
                // Skip malformed lines.
                return
            }
            hash := parts[0]
            date := strings.TrimSpace(parts[1])
            author := strings.TrimSpace(parts[2])
            subjectLine := parts[3]

            // Extract commit type from subject (if a colon exists).
            ctype := "unknown"
            if colonIdx := strings.Index(subjectLine, ":"); colonIdx != -1 {
                ctype = strings.TrimSpace(subjectLine[:colonIdx])
            }

            commitCh <- commitEntry{
                hash:    hash,
                date:    date,
                author:  author,
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
            // Include hash, date, author, and subject for each commit.
            buf.WriteString(fmt.Sprintf("- [%s] %s by %s: %s\n", e.hash, e.date, e.author, e.subject))
        }
        buf.WriteString("\n")
    }

    // Write the content to CHANGELOG.
    changelogPath := filepath.Join(".", "CHANGELOG")
    if err := os.WriteFile(changelogPath, buf.Bytes(), 0644); err != nil {
        return fmt.Errorf("failed to write CHANGELOG: %v", err)
    }
    fmt.Println("Changelog generated in", changelogPath)
    return nil
}

