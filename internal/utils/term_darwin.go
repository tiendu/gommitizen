//go:build darwin
// +build darwin

package utils

import (
    "os"
    "syscall"
    "unsafe"
)

// enableRawModeHelper enables raw mode on macOS.
func enableRawModeHelper() (*syscall.Termios, error) {
    fd := int(os.Stdin.Fd())
    termios := &syscall.Termios{}
    _, _, errno := syscall.Syscall6(
        syscall.SYS_IOCTL,
        uintptr(fd),
        uintptr(syscall.TIOCGETA),
        uintptr(unsafe.Pointer(termios)),
        0, 0, 0,
    )
    if errno != 0 {
        return nil, errno
    }
    raw := *termios
    raw.Lflag &^= syscall.ECHO | syscall.ICANON
    _, _, errno = syscall.Syscall6(
        syscall.SYS_IOCTL,
        uintptr(fd),
        uintptr(syscall.TIOCSETA),
        uintptr(unsafe.Pointer(&raw)),
        0, 0, 0,
    )
    if errno != 0 {
        return nil, errno
    }
    return termios, nil
}

// disableRawModeHelper restores terminal settings on macOS.
func disableRawModeHelper(orig *syscall.Termios) {
    fd := int(os.Stdin.Fd())
    syscall.Syscall6(
        syscall.SYS_IOCTL,
        uintptr(fd),
        uintptr(syscall.TIOCSETA),
        uintptr(unsafe.Pointer(orig)),
        0, 0, 0,
    )
}

// getTerminalSize returns terminal width and height on macOS.
func getTerminalSize() (width, height int, err error) {
    ws := &winsize{}
    _, _, errno := syscall.Syscall(
        syscall.SYS_IOCTL,
        os.Stdout.Fd(),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)),
    )
    if errno != 0 {
        return 0, 0, errno
    }
    return int(ws.cols), int(ws.rows), nil
}

