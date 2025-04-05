package cmd

import (
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    "gommitizen/internal/utils"
)

// VersionStr and RevisionStr can be set at build time during build.
var VersionStr = "0.3.0"
var RevisionStr = "unknown"

// HelpCommand displays help message for gommitizen CLI.
func HelpCommand() {
    fmt.Println(`Usage: gommitizen <command> [options]

Commands:
  install      Install this tool to ~/.local/bin/gommitizen and copy configs
               Options:
                 --path, -p  Specify custom install path (default: ~/.local/bin/gommitizen)

  reinstall    Uninstall and then reinstall the tool
               Options:
                 --path, -p  Specify custom install path (default: ~/.local/bin/gommitizen)

  uninstall    Remove the installed gommitizen
               Options:
                 --path, -p  Specify custom install path (default: ~/.local/bin/gommitizen)

  version      Print version information
  commit       Create a commit using the configured commitizen flow
  changelog    Generate a CHANGELOG.md from commit logs
  bump         Bump the version automatically
  lint         Lint commit messages
      Options for lint:
          --all, -a      Lint all commit messages in the repository
          --current, -c  Lint only the current (latest) commit message
          [commit-message]  Optionally, provide a commit message directly

  help         Display this help message`)
}

// InstallCommand installs this binary and its configs.
func InstallCommand(args []string) {
    installPath := parsePathFlag(args)

    appFilePath, err := os.Executable()
    if err != nil {
        fmt.Printf("Error determining executable path: %v\n", err)
        return
    }

    dest, err := installSubCmd(appFilePath, installPath)
    if err != nil {
        fmt.Printf("❌ Install failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Installed gommitizen to %s\n", dest)
    checkPathNotice(installPath)
}

// installSubCmd copies the binary and config to the install directory.
func installSubCmd(appFilePath, installPath string) (string, error) {
    configsPath := filepath.Join(installPath, "configs")
    if err := os.MkdirAll(configsPath, 0755); err != nil {
        return "", fmt.Errorf("failed to create install directory: %v", err)
    }

    destBinary := filepath.Join(installPath, "git-cz")
    if err := copyFile(appFilePath, destBinary, 0755); err != nil {
        return "", fmt.Errorf("failed to copy binary: %v", err)
    }

    srcConfig := filepath.Join("configs", "default.json")
    destConfig := filepath.Join(configsPath, "default.json")
    if err := copyFile(srcConfig, destConfig, 0644); err != nil {
        return "", fmt.Errorf("failed to copy config: %v", err)
    }

    return destBinary, nil
}

// UninstallCommand removes the installed gommitizen directory.
func UninstallCommand(args []string) {
    installPath := parsePathFlag(args)

    if _, err := os.Stat(installPath); os.IsNotExist(err) {
        fmt.Printf("Nothing to uninstall at %s\n", installPath)
        return
    }

    if err := os.RemoveAll(installPath); err != nil {
        fmt.Printf("❌ Uninstall failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Uninstalled gommitizen from %s\n", installPath)
}

// ReinstallCommand uninstalls and reinstalls the tool.
func ReinstallCommand(args []string) {
    installPath := parsePathFlag(args)

    UninstallCommand(args)
    InstallCommand(args)

    fmt.Printf("✅ Reinstallation completed at %s\n", installPath)
}

// VersionCommand prints version information.
func VersionCommand() {
    fmt.Printf("gommitizen version %s, build revision %s\n", VersionStr, RevisionStr)
}

// LintOptions holds the options for the lint command.
type LintOptions struct {
    All     bool
    Current bool
    Message string
}

// ParseLintOptions parses the lint command flags and returns a LintOptions struct.
func ParseLintOptions(args []string) (LintOptions, error) {
    lf := flag.NewFlagSet("lint", flag.ExitOnError)
    allLong := lf.Bool("all", false, "Lint all commit messages")
    allShort := lf.Bool("a", false, "Lint all commit messages (short)")
    currentLong := lf.Bool("current", false, "Lint the current commit message")
    currentShort := lf.Bool("c", false, "Lint the current commit message (short)")
    lf.Parse(args)

    opts := LintOptions{
        All:     *allLong || *allShort,
        Current: *currentLong || *currentShort,
    }

    if lf.NArg() > 0 {
        opts.Message = lf.Arg(0)
    }

    return opts, nil
}

// Helpers

// parsePathFlag parses the --path or -p flag from args.
func parsePathFlag(args []string) string {
    fs := flag.NewFlagSet("path", flag.ExitOnError)
    path := fs.String("path", defaultInstallPath(), "Specify custom install path")
    pathShort := fs.String("p", defaultInstallPath(), "Specify custom install path (shorthand)")
    _ = fs.Parse(args)

    if *pathShort != defaultInstallPath() {
        return *pathShort
    }
    return *path
}

// defaultInstallPath returns the default installation path (~/.local/bin/gommitizen)
func defaultInstallPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./gommitizen"
    }
    return filepath.Join(homeDir, ".local", "bin", "gommitizen")
}

// copyFile copies a file from src to dst with given permissions.
func copyFile(srcPath, destPath string, mode os.FileMode) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, srcFile)
    return err
}

// checkPathNotice advises the user to add install path to PATH if it's not already.
func checkPathNotice(installPath string) {
    if isInPath(installPath) {
        return
    }

    shell := detectShell()

    fmt.Println()
    fmt.Println(utils.Bold(utils.Color("⚠️  PATH Notice:", "yellow")))
    fmt.Printf("Your install path %s is %s in your PATH.\n",
        utils.Color(installPath, "cyan"),
        utils.Color("not", "red"),
    )

    fmt.Println(utils.Color("\nTo use gommitizen from anywhere, add it to your PATH:", "magenta"))

    switch shell {
    case "zsh":
        fmt.Printf("\n%s\n", utils.Color(
            fmt.Sprintf("echo 'export PATH=\"%s:$PATH\"' >> ~/.zshrc && source ~/.zshrc", installPath),
            "green",
        ))
    case "bash":
        fmt.Printf("\n%s\n", utils.Color(
            fmt.Sprintf("echo 'export PATH=\"%s:$PATH\"' >> ~/.bashrc && source ~/.bashrc", installPath),
            "green",
        ))
    default:
        fmt.Println(utils.Color(
            fmt.Sprintf("\nYour shell could not be detected. Please manually add %s to your PATH.", installPath),
            "yellow",
        ))
    }

    fmt.Println(utils.Color("\nAfter adding, restart your terminal or run the above command to apply changes immediately.", "yellow"))
    fmt.Println()
}

// detectShell tries to detect the user's shell from environment.
func detectShell() string {
    shellEnv := os.Getenv("SHELL")
    if strings.Contains(shellEnv, "zsh") {
        return "zsh"
    }
    if strings.Contains(shellEnv, "bash") {
        return "bash"
    }
    return "unknown"
}

// isInPath checks if the given directory is in the PATH environment variable.
func isInPath(dir string) bool {
    pathEnv := os.Getenv("PATH")
    for _, p := range filepath.SplitList(pathEnv) {
        if p == dir {
            return true
        }
    }
    return false
}

