package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func resetGlobals() {
	isScreen = false
	isTmux = false
	isZellij = false
}

func TestOscSequence(t *testing.T) {
	defer resetGlobals()

	t.Run("default OSC 0", func(t *testing.T) {
		resetGlobals()
		open, close := oscSequence("0")
		expectedOpen := "\x1b]0;"
		expectedClose := "\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("default OSC 1", func(t *testing.T) {
		resetGlobals()
		open, close := oscSequence("1")
		expectedOpen := "\x1b]1;"
		expectedClose := "\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("default OSC 2", func(t *testing.T) {
		resetGlobals()
		open, close := oscSequence("2")
		expectedOpen := "\x1b]2;"
		expectedClose := "\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("default OSC 52", func(t *testing.T) {
		resetGlobals()
		open, close := oscSequence("52;c")
		expectedOpen := "\x1b]52;c;"
		expectedClose := "\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("screen OSC 0", func(t *testing.T) {
		resetGlobals()
		isScreen = true
		open, close := oscSequence("0")
		expectedOpen := "\x1bP\x1b]0;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("screen OSC 1", func(t *testing.T) {
		resetGlobals()
		isScreen = true
		open, close := oscSequence("1")
		expectedOpen := "\x1bP\x1b]1;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("screen OSC 2", func(t *testing.T) {
		resetGlobals()
		isScreen = true
		open, close := oscSequence("2")
		expectedOpen := "\x1bP\x1b]2;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("tmux OSC 0", func(t *testing.T) {
		resetGlobals()
		isTmux = true
		open, close := oscSequence("0")
		expectedOpen := "\x1bPtmux;\x1b\x1b]0;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("tmux OSC 1", func(t *testing.T) {
		resetGlobals()
		isTmux = true
		open, close := oscSequence("1")
		expectedOpen := "\x1bPtmux;\x1b\x1b]1;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})

	t.Run("tmux OSC 2", func(t *testing.T) {
		resetGlobals()
		isTmux = true
		open, close := oscSequence("2")
		expectedOpen := "\x1bPtmux;\x1b\x1b]2;"
		expectedClose := "\x1b\\\x1b\\"
		if open != expectedOpen {
			t.Errorf("expected open %q, got %q", expectedOpen, open)
		}
		if close != expectedClose {
			t.Errorf("expected close %q, got %q", expectedClose, close)
		}
	})
}

func TestWriteTitle(t *testing.T) {
	defer resetGlobals()

	for _, oscCode := range []string{"0", "1", "2"} {
		t.Run("normal terminal OSC "+oscCode+" single word", func(t *testing.T) {
			resetGlobals()
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "test")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1b]%s;test\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("normal terminal OSC "+oscCode+" multiple words", func(t *testing.T) {
			resetGlobals()
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "my terminal title")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1b]%s;my terminal title\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("normal terminal OSC "+oscCode+" empty string", func(t *testing.T) {
			resetGlobals()
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1b]%s;\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("tmux terminal OSC "+oscCode, func(t *testing.T) {
			resetGlobals()
			isTmux = true
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "my terminal title")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1bPtmux;\x1b\x1b]%s;my terminal title\x1b\\\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("tmux terminal OSC "+oscCode+" with internal ESC", func(t *testing.T) {
			resetGlobals()
			isTmux = true
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "esc\x1btest")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1bPtmux;\x1b\x1b]%s;esc\x1b\x1btest\x1b\\\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("screen terminal OSC "+oscCode+" short title", func(t *testing.T) {
			resetGlobals()
			isScreen = true
			var buf bytes.Buffer
			err := writeTitle(&buf, oscCode, "short title")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1bP\x1b]%s;short title\x1b\\\x1b\\", oscCode)
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})

		t.Run("screen terminal OSC "+oscCode+" long title chunking", func(t *testing.T) {
			resetGlobals()
			isScreen = true
			var buf bytes.Buffer
			longTitle := strings.Repeat("a", 100)
			err := writeTitle(&buf, oscCode, longTitle)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := fmt.Sprintf("\x1bP\x1b]%s;%s\x1b\\\x1bP%s\x1b\\\x1b\\", oscCode, strings.Repeat("a", 76), strings.Repeat("a", 24))
			if buf.String() != expected {
				t.Errorf("expected %q, got %q", expected, buf.String())
			}
		})
	}
}

func TestTitleCmdArgsValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  *cobra.Command
	}{
		{"title-0", title0Cmd},
		{"title-1", title1Cmd},
		{"title-2", title2Cmd},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.cmd.Args(tc.cmd, []string{}); err == nil {
				t.Error("expected error for 0 arguments, got nil")
			}
			if err := tc.cmd.Args(tc.cmd, []string{"hello"}); err != nil {
				t.Errorf("expected nil for 1 argument, got: %v", err)
			}
			if err := tc.cmd.Args(tc.cmd, []string{"hello", "world"}); err != nil {
				t.Errorf("expected nil for 2 arguments, got: %v", err)
			}
			if err := tc.cmd.Args(tc.cmd, []string{""}); err != nil {
				t.Errorf("expected nil for empty string argument, got: %v", err)
			}
		})
	}
}

func TestCLIHelpAndValidation(t *testing.T) {
	tempDir := t.TempDir()
	binPath := filepath.Join(tempDir, "osc")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v, output: %s", err, string(out))
	}

	t.Run("root help lists title commands", func(t *testing.T) {
		cmd := exec.Command(binPath, "--help")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("osc --help failed: %v, output: %s", err, string(out))
		}
		outStr := string(out)
		for _, name := range []string{"title-2", "title-0", "title-1"} {
			if !strings.Contains(outStr, name) {
				t.Errorf("osc --help output missing %q: %s", name, outStr)
			}
		}
	})

	t.Run("title help shows usage", func(t *testing.T) {
		for _, tc := range []struct {
			cmdName string
			oscCode string
		}{
			{"title-0", "0"},
			{"title-1", "1"},
			{"title-2", "2"},
		}{
			cmd := exec.Command(binPath, tc.cmdName, "--help")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("osc %s --help failed: %v, output: %s", tc.cmdName, err, string(out))
			}
			outStr := string(out)
			expectedUsage := fmt.Sprintf("osc %s <str> [str...]", tc.cmdName)
			if !strings.Contains(outStr, expectedUsage) {
				t.Errorf("osc %s --help missing expected usage %q: %s", tc.cmdName, expectedUsage, outStr)
			}
			expectedOSC := fmt.Sprintf("OSC %s", tc.oscCode)
			if !strings.Contains(outStr, expectedOSC) {
				t.Errorf("osc %s --help missing %q: %s", tc.cmdName, expectedOSC, outStr)
			}
		}
	})

	t.Run("title subcommands require at least 1 arg", func(t *testing.T) {
		for _, cmdName := range []string{"title-2", "title-0", "title-1"} {
			cmd := exec.Command(binPath, cmdName)
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("expected osc %s with no args to fail, but succeeded", cmdName)
			}
			outStr := string(out)
			if !strings.Contains(outStr, "requires at least 1 arg(s), only received 0") {
				t.Errorf("unexpected error output for %s: %s", cmdName, outStr)
			}
		}
	})
}
