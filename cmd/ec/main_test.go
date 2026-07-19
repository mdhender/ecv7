package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mdhender/ecv7"
	"github.com/peterbourgon/ff/v4"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "core", args: []string{"version"}, want: ecv7.Version().Core()},
		{name: "build", args: []string{"version", "--build"}, want: ecv7.Version().Short()},
		{name: "long", args: []string{"version", "--long"}, want: ecv7.Version().String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := captureStdout(t, func() error {
				return run(t.Context(), tt.args, &bytes.Buffer{})
			})
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			if got := strings.TrimSpace(output); got != tt.want {
				t.Fatalf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVersionFlagsAreMutuallyExclusive(t *testing.T) {
	_, err := captureStdout(t, func() error {
		return run(t.Context(), []string{"version", "--build", "--long"}, &bytes.Buffer{})
	})
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("run error = %v, want mutually exclusive error", err)
	}
}

func TestRootFlagsAreInherited(t *testing.T) {
	for _, args := range [][]string{
		{"--quiet", "--path", "games/example/db", "version"},
		{"version", "--quiet", "--path", "games/example/db"},
	} {
		output, err := captureStdout(t, func() error {
			return run(t.Context(), args, &bytes.Buffer{})
		})
		if err != nil {
			t.Fatalf("run(%v): %v", args, err)
		}
		if got := strings.TrimSpace(output); got != ecv7.Version().Core() {
			t.Fatalf("run(%v) output = %q, want %q", args, got, ecv7.Version().Core())
		}
	}
}

func TestHelp(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"version", "--help"}} {
		var stderr bytes.Buffer
		err := run(t.Context(), args, &stderr)
		if !errors.Is(err, ff.ErrHelp) {
			t.Fatalf("run(%v) error = %v, want ErrHelp", args, err)
		}
		help := stderr.String()
		for _, flag := range []string{"--path", "--quiet"} {
			if !strings.Contains(help, flag) {
				t.Errorf("run(%v) help = %q, want %s", args, help, flag)
			}
		}
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = original })

	runErr := fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = original
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(output), runErr
}
