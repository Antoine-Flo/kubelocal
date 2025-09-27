package ui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// printChoicePrompt displays the question with default choice
func printChoicePrompt(question, defaultChoice string) {
	fmt.Printf("%s%s%s%s%s%s", ColorWhite, question, ColorReset, ColorGray, defaultChoice, ColorReset)
}

// readRawChar reads a single character from stdin
func readRawChar() (byte, error) {
	var b [1]byte
	_, err := os.Stdin.Read(b[:])
	return b[0], err
}

// isEnterKey checks if the character is Enter key
func isEnterKey(char byte) bool {
	return char == 13 || char == 10
}

// printSelectedChoice displays the final choice
func printSelectedChoice(question, choice string) {
	fmt.Printf("\r\033[K%s%s%s%s%s\n", ColorWhite, question, ColorReset, ColorGreen, choice)
	fmt.Print("\r")
}

// setRawMode sets terminal to raw mode
func setRawMode() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

// restoreTerminal restores terminal to original state
func restoreTerminal(oldState *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), oldState)
}

// AskChoiceWithDefault prompts user for choice with default option
func AskChoiceWithDefault(question string, options map[string]string, defaultChoice string) (string, error) {
	printChoicePrompt(question, defaultChoice)

	oldState, err := setRawMode()
	if err != nil {
		return defaultChoice, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer restoreTerminal(oldState)

	return readUserChoice(question, options, defaultChoice, oldState)
}

// readUserChoice reads and processes user input
func readUserChoice(question string, options map[string]string, defaultChoice string, oldState *term.State) (string, error) {
	for {
		char, err := readRawChar()
		if err != nil {
			continue
		}

		if choice := processInput(char, options, defaultChoice, oldState); choice != "" {
			printSelectedChoice(question, choice)
			return choice, nil
		}
	}
}

// processInput processes a character input and returns choice if valid
func processInput(char byte, options map[string]string, defaultChoice string, oldState *term.State) string {
	// Ctrl+C = ASCII 3
	if char == 3 {
		restoreTerminal(oldState)
		os.Exit(0)
	}

	if isEnterKey(char) {
		return defaultChoice
	}

	input := strings.ToLower(string(char))
	if _, exists := options[input]; exists {
		return input
	}

	return ""
}
