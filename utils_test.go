package main

import (
	"strings"
	"testing"
)

func TestWrapSingleLine(t *testing.T) {
	line := "one two three four"
	got := wrapSingleLine(line, 8)
	want := "one two\nthree\nfour"
	if got != want {
		t.Errorf("wrapSingleLine() = %q, want %q", got, want)
	}

	if got := wrapSingleLine(line, 50); got != line {
		t.Errorf("wrapSingleLine() = %q, want %q", got, line)
	}
}

func TestSplitIntoLines(t *testing.T) {
	input := "a\nb\r\nc"
	want := []string{"a", "b", "c"}
	got := splitIntoLines(input)
	if len(got) != len(want) {
		t.Fatalf("splitIntoLines length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("splitIntoLines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestWrap(t *testing.T) {
	text := "hello world\n\nfoo bar baz"
	want := "hello\nworld\n\nfoo\nbar\nbaz"
	if got := wrap(text, 6); got != want {
		t.Errorf("wrap() = %q, want %q", got, want)
	}
}

func TestExampleCommands(t *testing.T) {
	examples := []example{{"-h", "help"}, {"--version", "print version"}}
	output := exampleCommands("cmd", examples)
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	expected := []string{
		"  cmd -h            help",
		"  cmd --version     print version",
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line %d = %q, want %q", i, line, expected[i])
		}
	}
}
