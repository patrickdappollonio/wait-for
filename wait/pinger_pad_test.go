package wait

import (
	"io"
	"os"
	"testing"
)

func TestAppPad(t *testing.T) {
	a := &App{padding: 5}
	if got := a.pad("a"); got != "a    " {
		t.Errorf("pad() = %q", got)
	}
}

func TestPrintOnVerbose(t *testing.T) {
	a := &App{Verbose: true}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	a.printOnVerbose("msg")
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if string(out) != "msg\n" {
		t.Errorf("unexpected output: %q", string(out))
	}

	// now ensure nothing printed when verbose false
	r, w, _ = os.Pipe()
	os.Stdout = w
	a.Verbose = false
	a.printOnVerbose("nope")
	w.Close()
	os.Stdout = old
	out, _ = io.ReadAll(r)
	if len(out) != 0 {
		t.Errorf("expected no output, got %q", string(out))
	}
}
