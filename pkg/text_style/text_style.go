package text_style

import "runtime"

// Color constants
const (
	Reset   = "\033[0m"
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	Yellow  = "\033[0;33m"
	Blue    = "\033[0;34m"
	Magenta = "\033[0;35m"
	Cyan    = "\033[0;36m"
	White   = "\033[0;37m"
)

func ColorText(color, text string) string {
	// colorText wraps a string with the given color if supported
	if runtime.GOOS == "windows" && !isWindowsTerminalWithANSI() {
		// If running on older Windows versions without ANSI support, return plain text
		return text
	}
	return color + text + Reset
}

func BoldText(text string) string {
	// boldText wraps a string with ANSI bold codes if supported
	if runtime.GOOS == "windows" && !isWindowsTerminalWithANSI() {
		// If running on older Windows versions without ANSI support, return plain text
		return text
	}
	return "\033[1m" + text + Reset
}

func isWindowsTerminalWithANSI() bool {
	// Detect if the terminal supports ANSI on modern Windows (this is a simplified check)
	return runtime.GOOS == "windows" && (hasANSISupportInWindows10() || runningInWindowsTerminal())
}

func hasANSISupportInWindows10() bool {
	// Mock function: Implement actual detection of Windows 10 with ANSI support if necessary
	return true // Assuming Windows 10 for now
}

func runningInWindowsTerminal() bool {
	// Mock function: Detect if running in Windows Terminal or another ANSI-capable terminal
	return true // Assuming modern Windows Terminal for now
}
