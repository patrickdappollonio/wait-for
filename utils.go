package main

import (
	"fmt"
	"strings"
)

// wrap wraps the input text to the specified width without splitting words.
// It preserves existing line breaks.
// If width is less than or equal to 0, it defaults to 80.
func wrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	var builder strings.Builder

	// Split the text into lines based on newline characters.
	// Handles both Unix (\n) and Windows (\r\n) line endings.
	lines := splitIntoLines(text)

	for i, line := range lines {
		// For lines after the first, prepend a newline to preserve line breaks.
		if i > 0 {
			builder.WriteByte('\n')
		}

		// If the line is empty, preserve the empty line.
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Wrap the individual line and write to the builder.
		wrapped := wrapSingleLine(line, width)
		builder.WriteString(wrapped)
	}

	return builder.String()
}

// splitIntoLines splits the input text into lines, handling both \n and \r\n.
func splitIntoLines(text string) []string {
	// Normalize line endings to \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Split by \n
	return strings.Split(text, "\n")
}

// wrapSingleLine wraps a single line of text to the specified width without splitting words.
func wrapSingleLine(line string, width int) string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return ""
	}

	var wrapped strings.Builder
	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding the next word exceeds the width
		if len(currentLine)+1+len(word) > width {
			// Write the current line to the builder
			wrapped.WriteString(currentLine)
			wrapped.WriteByte('\n')
			// Start a new line with the current word
			currentLine = word
		} else {
			// Add the word to the current line
			currentLine += " " + word
		}
	}

	// Append the last line
	wrapped.WriteString(currentLine)

	return wrapped.String()
}

type example struct {
	command string
	helper  string
}

func exampleCommands(cmdname string, examples []example) string {
	padding := 0
	for _, v := range examples {
		if len(v.command) > padding {
			padding = len(v.command)
		}
	}

	var sb strings.Builder

	for i, v := range examples {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(
			&sb,
			"  %s %s %s %s",
			cmdname, v.command, strings.Repeat(" ", padding-len(v.command)+3), v.helper,
		)
	}

	return sb.String()
}
