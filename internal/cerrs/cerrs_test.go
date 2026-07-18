package cerrs

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	const err Error = "constant error"

	if got := err.Error(); got != "constant error" {
		t.Fatalf("Error() = %q, want %q", got, "constant error")
	}
	if !errors.Is(err, Error("constant error")) {
		t.Fatal("errors.Is did not match equal constant errors")
	}
}
