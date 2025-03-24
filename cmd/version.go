package cmd

import (
    "os"
    "strconv"
    "strings"
    "fmt"
)

// BumpVersion reads the current version from the VERSION file and dynamically bumps the version:
// - If version is in the format X.Y, it increments Y (minor version).
// - If version is in the format X.Y.Z, it increments Z (patch version).
func BumpVersion() (string, error) {
    // Read the current version from the VERSION file.
    data, err := os.ReadFile("VERSION")
    if err != nil {
        return "", fmt.Errorf("failed to read VERSION file: %v", err)
    }
    versionStr := strings.TrimSpace(string(data))
    parts := strings.Split(versionStr, ".")
    var newVersion string

    if len(parts) == 2 {
        // Format: X.Y -> increment the minor version.
        major, err := strconv.Atoi(parts[0])
        if err != nil {
            return "", fmt.Errorf("failed to parse major version: %v", err)
        }
        minor, err := strconv.Atoi(parts[1])
        if err != nil {
            return "", fmt.Errorf("failed to parse minor version: %v", err)
        }
        minor++
        newVersion = fmt.Sprintf("%d.%d", major, minor)
    } else if len(parts) == 3 {
        // Format: X.Y.Z -> increment the patch version.
        major, err := strconv.Atoi(parts[0])
        if err != nil {
            return "", fmt.Errorf("failed to parse major version: %v", err)
        }
        minor, err := strconv.Atoi(parts[1])
        if err != nil {
            return "", fmt.Errorf("failed to parse minor version: %v", err)
        }
        patch, err := strconv.Atoi(parts[2])
        if err != nil {
            return "", fmt.Errorf("failed to parse patch version: %v", err)
        }
        patch++
        newVersion = fmt.Sprintf("%d.%d.%d", major, minor, patch)
    } else {
        return "", fmt.Errorf("version format invalid, expected X.Y or X.Y.Z")
    }

    // Write the new version back to the VERSION file.
    if err := os.WriteFile("VERSION", []byte(newVersion), 0644); err != nil {
        return "", fmt.Errorf("failed to write VERSION file: %v", err)
    }

    fmt.Printf("Version bumped from %s to %s\n", versionStr, newVersion)

    return newVersion, nil
}

