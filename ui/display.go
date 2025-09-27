package ui

import (
	"fmt"
	"strings"
)

// PrintHeader displays a colored header with underline
func PrintHeader(title string) {
	fmt.Println(strings.Repeat("=", len(title)))
	fmt.Printf("%s%s%s\n", ColorCyan, title, ColorReset)
	fmt.Println(strings.Repeat("=", len(title)))
	fmt.Println()
}

// PrintSuccess displays a success message with checkmark
func PrintSuccess(message string) {
	fmt.Printf("%s✅ %s%s\n", ColorGreen, message, ColorReset)
}

// PrintStep displays a step counter with message
func PrintStep(step, total int, message string) {
	fmt.Printf("%s[%d/%d]%s %s... ", ColorBlue, step, total, ColorReset, message)
}
