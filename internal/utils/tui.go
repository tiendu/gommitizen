package utils

import (
    "fmt"
    "os"
    "os/signal"
    "regexp"
    "syscall"
    "unsafe"
    "time"
    "math"
)

// ========================
// Terminal Utilities
// ========================

// winsize holds terminal size (rows/cols)
type winsize struct {
    rows uint16
    cols uint16
    x    uint16
    y    uint16
}

// getTerminalHeight returns the terminal's height (rows).
func getTerminalHeight() (int, error) {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdout),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)),
    )
    if int(retCode) == -1 {
        return 0, errno
    }
    return int(ws.rows), nil
}

// isTerminal checks if stdin is a terminal.
func isTerminal() bool {
    fi, err := os.Stdin.Stat()
    if err != nil {
        return false
    }
    return (fi.Mode() & os.ModeCharDevice) != 0
}

// ========================
// Terminal UI Struct
// ========================

type TerminalUI struct {
    originalTermios   *syscall.Termios // Original terminal state
    selectedIndex     int              // Current selected index
    menuOptions       []string         // List of menu options
    prevRenderBuffer  []string         // Cached previous render (for diffing)
    maxVisibleOptions int              // Max visible options at once
    firstVisibleIndex int              // Start index for visible window
    maxOptionWidth    int              // Max width for wrapped lines
    cleanupOnce       bool             // Ensure cleanup runs once
    anchorRow         int              // Starting row for rendering
    terminalHeight    int              // Terminal height (rows)
}

// ========================
// Setup / Teardown Helpers
// ========================

func (t *TerminalUI) setup() error {
    if !isTerminal() {
        return fmt.Errorf("stdin is not a terminal")
    }

    // Hide the cursor
    fmt.Print("\033[?25l")

    // Enable raw mode
    termios, err := enableRawModeHelper()
    if err != nil {
        return err
    }
    t.originalTermios = termios

    // Get terminal height
    height, err := getTerminalHeight()
    if err != nil {
        return err
    }
    t.terminalHeight = height

    // Get current cursor position (to align TUI below existing output)
    row, _, err := t.getCursorPosition()
    if err != nil {
        return err
    }

    // Estimate menu height (header + options with wrapping space)
    estimatedMenuHeight := 1 + t.maxVisibleOptions*3
    proposedAnchorRow := row + 1

    // If menu would exceed terminal height, scroll terminal to make room
    if proposedAnchorRow+estimatedMenuHeight > t.terminalHeight {
        linesToScroll := proposedAnchorRow + estimatedMenuHeight - t.terminalHeight
        for i := 0; i < linesToScroll; i++ {
            fmt.Println()
        }
        proposedAnchorRow -= linesToScroll
        if proposedAnchorRow < 1 {
            proposedAnchorRow = 1
        }
    }

    t.anchorRow = proposedAnchorRow

    // Handle Ctrl+C and termination signals
    t.handleSignals()

    return nil
}

func (t *TerminalUI) cleanup() {
    if t.cleanupOnce {
        return
    }
    t.cleanupOnce = true

    t.disableRawMode()

    // Show cursor back
    fmt.Print("\033[?25h")
}

func (t *TerminalUI) disableRawMode() {
    if t.originalTermios != nil {
        disableRawModeHelper(t.originalTermios)
    }
}

// ========================
// Terminal Position Helpers
// ========================

// Get current cursor position by issuing DSR request
func (t *TerminalUI) getCursorPosition() (int, int, error) {
    fmt.Print("\033[6n") // Request: "Report Cursor Position"
    var buf []byte
    tmp := make([]byte, 1)
    for {
        _, err := os.Stdin.Read(tmp)
        if err != nil {
            return 0, 0, err
        }
        buf = append(buf, tmp[0])
        if tmp[0] == 'R' {
            break
        }
    }
    var row, col int
    _, err := fmt.Sscanf(string(buf[2:]), "%d;%d", &row, &col)
    return row, col, err
}

// Move cursor to row, col
func (t *TerminalUI) moveCursor(row, col int) {
    fmt.Printf("\033[%d;%dH", row, col)
}

// ========================
// Signal Handling
// ========================

func (t *TerminalUI) handleSignals() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        <-c
        t.cleanup()
        os.Exit(1)
    }()
}

// ========================
// Input Handling
// ========================

// readKey reads a single key press.
func (t *TerminalUI) readKey() (rune, error) {
    var buf [3]byte
    n, err := os.Stdin.Read(buf[:])
    if err != nil {
        return 0, err
    }
    if n == 3 && buf[0] == 27 && buf[1] == 91 {
        switch buf[2] {
        case 'A':
            return '↑', nil
        case 'B':
            return '↓', nil
        }
    }
    return rune(buf[0]), nil
}

// ========================
// Rendering
// ========================

func (t *TerminalUI) renderMenu() {
    t.moveCursor(t.anchorRow, 1)

    currentRender := []string{
        fmt.Sprintf("Use ↑ (k) ↓ (j) to move, Enter to select: (Item %d of %d)", t.selectedIndex+1, len(t.menuOptions)),
    }

    // Maintain visible window
    if t.selectedIndex < t.firstVisibleIndex {
        t.firstVisibleIndex = t.selectedIndex
    } else if t.selectedIndex >= t.firstVisibleIndex+t.maxVisibleOptions {
        t.firstVisibleIndex = t.selectedIndex - t.maxVisibleOptions + 1
    }

    end := t.firstVisibleIndex + t.maxVisibleOptions
    if end > len(t.menuOptions) {
        end = len(t.menuOptions)
    }

    // Render options with wrapping
    for i := t.firstVisibleIndex; i < end; i++ {
        wrappedLines := breakLines(t.menuOptions[i], t.maxOptionWidth)
        if len(wrappedLines) == 0 {
            wrappedLines = []string{""}
        }

        if i == t.selectedIndex {
            pointer := Color("❯", "green")

            // First line: pointer + underline + highlight
            firstLine := Highlight(Underline(StripANSI(wrappedLines[0])), "white", "green")
            currentRender = append(currentRender, fmt.Sprintf("%s %s", pointer, firstLine))

            // Following lines: indent + underline
            for _, line := range wrappedLines[1:] {
                underlined := Highlight(Underline(StripANSI(line)), "white", "green")
                currentRender = append(currentRender, fmt.Sprintf("    %s", underlined))
            }

        } else {
            // Non-selected: plain clean lines
            cleanLine := StripANSI(wrappedLines[0])
            currentRender = append(currentRender, "  " + cleanLine)

            for _, line := range wrappedLines[1:] {
                cleanLine := StripANSI(line)
                currentRender = append(currentRender, "    " + cleanLine)
            }
        }
    }

    // Diff rendering: update only changed lines
    for i, line := range currentRender {
        var prev string
        if i < len(t.prevRenderBuffer) {
            prev = t.prevRenderBuffer[i]
        }
        if line != prev {
            t.moveCursor(t.anchorRow+i, 1)
            fmt.Print("\033[2K")
            fmt.Println(line)
        }
    }

    // Clean up leftover lines from previous render
    for i := len(currentRender); i < len(t.prevRenderBuffer); i++ {
        t.moveCursor(t.anchorRow+i, 1)
        fmt.Print("\033[2K\r\n")
    }

    t.prevRenderBuffer = currentRender
}

func (t *TerminalUI) clearRenderedArea() {
    t.moveCursor(t.anchorRow, 1)
    for i := 0; i < len(t.prevRenderBuffer); i++ {
        t.moveCursor(t.anchorRow+i, 1)
        fmt.Print("\033[2K")
    }
}

// ========================
// Main Event Loop
// ========================

func (t *TerminalUI) Run() (int, string, error) {
    defer t.cleanup()

    if err := t.setup(); err != nil {
        return -1, "", err
    }

    t.renderMenu()

    for {
        key, err := t.readKey()
        if err != nil {
            return -1, "", fmt.Errorf("input read error: %v", err)
        }

        targetIndex := t.selectedIndex

        switch key {
        case '↑', 'k':
            targetIndex = t.selectedIndex - 1
            if targetIndex < 0 {
                targetIndex = len(t.menuOptions) - 1
            }

        case '↓', 'j':
            targetIndex = t.selectedIndex + 1
            if targetIndex >= len(t.menuOptions) {
                targetIndex = 0
            }

        case '\r', '\n':
            t.moveCursor(t.anchorRow+len(t.prevRenderBuffer)+1, 1)
            fmt.Println("Selected:", t.menuOptions[t.selectedIndex])
            return t.selectedIndex, t.menuOptions[t.selectedIndex], nil

        case 3: // Ctrl+C
            return -1, "", fmt.Errorf("cancelled by user")
        }

        // Smooth step-by-step animation
        if targetIndex != t.selectedIndex {
            step := 1
            if targetIndex < t.selectedIndex {
                step = -1
            }

            steps := abs(targetIndex - t.selectedIndex)
            for i := 1; i <= steps; i++ {
                t.selectedIndex += step
                t.renderMenu()

                // Easing effect
                delay := time.Duration(5 + int(math.Sqrt(float64(steps-i)) * 10)) * time.Millisecond
                time.Sleep(delay)
            }
        }

        t.renderMenu()
    }
}

// ========================
// Text Wrapping Helpers
// ========================

func splitWord(word string, maxLen int) []string {
    if len(word) <= maxLen {
        return []string{word}
    }
    var parts []string
    for len(word) > maxLen-1 {
        parts = append(parts, word[:maxLen-1]+"-")
        word = word[maxLen-1:]
    }
    parts = append(parts, word)
    return parts
}

func breakLines(text string, n int) []string {
    re := regexp.MustCompile(`\S+|[.,!?;()'"\[\]{}]`)
    words := re.FindAllString(text, -1)
    var lines []string
    var currentLine string
    for _, word := range words {
        if len(word) > n {
            chunks := splitWord(word, n)
            for _, chunk := range chunks {
                if len(currentLine)+len(chunk)+1 > n && currentLine != "" {
                    lines = append(lines, currentLine)
                    currentLine = chunk
                } else if currentLine != "" {
                    currentLine += " " + chunk
                } else {
                    currentLine = chunk
                }
            }
        } else if len(currentLine)+len(word)+1 > n && currentLine != "" {
            lines = append(lines, currentLine)
            currentLine = word
        } else if currentLine != "" {
            currentLine += " " + word
        } else {
            currentLine = word
        }
    }
    if currentLine != "" {
        lines = append(lines, currentLine)
    }
    return lines
}

func NewSelector(options []string, visible int, width int) *TerminalUI {
    return &TerminalUI{
        menuOptions:       options,
        maxVisibleOptions: visible,
        maxOptionWidth:    width,
    }
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
