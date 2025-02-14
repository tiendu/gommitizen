package internal

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// VersionStr and RevisionStr can be set at build time.
var VersionStr = "0.1.0"
var RevisionStr = "unknown"

// InstallCommand installs this binary as a Git subcommand (git-cz).
func InstallCommand() {
    appFilePath, err := os.Executable()
    if err != nil {
        fmt.Printf("Error determining executable path: %v\n", err)
        return
    }
    dest, err := installSubCmd(appFilePath, "cz")
    if err != nil {
        fmt.Printf("Install commitizen failed, err=%v\n", err)
    } else {
        fmt.Printf("Installed commitizen to %s\n", dest)
    }
}

// installSubCmd copies the current binary to Git's exec-path with the name "git-<subCmd>".
func installSubCmd(appFilePath, subCmd string) (string, error) {
    out, err := exec.Command("git", "--exec-path").Output()
    if err != nil {
        return "", err
    }
    execPath := strings.TrimSpace(string(out))
    destPath := filepath.Join(execPath, "git-"+subCmd)

    // Ensure the exec directory exists.
    if _, err := os.Stat(execPath); os.IsNotExist(err) {
        if err := os.MkdirAll(execPath, 0755); err != nil {
            return "", fmt.Errorf("failed to create directory %s: %v", execPath, err)
        }
    }

    srcFile, err := os.Open(appFilePath)
    if err != nil {
        return "", err
    }
    defer srcFile.Close()

    destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0755)
    if err != nil {
        return "", err
    }
    defer destFile.Close()

    if _, err := io.Copy(destFile, srcFile); err != nil {
        return "", err
    }
    return destPath, nil
}

// UninstallCommand removes the installed Git subcommand (git-cz).
func UninstallCommand() {
    out, err := exec.Command("git", "--exec-path").Output()
    if err != nil {
        fmt.Printf("Error obtaining git exec path: %v\n", err)
        return
    }
    execPath := strings.TrimSpace(string(out))
    destPath := filepath.Join(execPath, "git-cz")

    if err := os.Remove(destPath); err != nil {
        fmt.Printf("Uninstall failed: %v\n", err)
        return
    }
    fmt.Printf("Uninstalled commitizen from %s\n", destPath)
}

// ReinstallCommand uninstalls and then reinstalls the tool.
func ReinstallCommand() {
    UninstallCommand()
    InstallCommand()
}

// VersionCommand prints version information.
func VersionCommand() {
    fmt.Printf("Commitizen version %s, build revision %s\n", VersionStr, RevisionStr)
}

