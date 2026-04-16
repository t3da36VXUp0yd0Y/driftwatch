package redact_test

import (
	"testing"

	"github.com/driftwatch/internal/redact"
)

func TestSanitise_NoSensitiveData(t *testing.T) {
	r := redact.New()
	input := "registry.example.com/myimage:latest"
	if got := r.Sanitise(input); got != input {
		t.Errorf("expected %q unchanged, got %q", input, got)
	}
}

func TestSanitise_BasicAuthURL(t *testing.T) {
	r := redact.New()
	input := "https://user:s3cr3t@registry.example.com/image:tag"
	got := r.Sanitise(input)
	if got != "https://***@registry.example.com/image:tag" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestSanitise_ExtraSecret(t *testing.T) {
	r := redact.New("mysecrettoken")
	input := "image built with mysecrettoken at build time"
	got := r.Sanitise(input)
	expected := "image built with *** at build time"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSanitise_EmptySecretIgnored(t *testing.T) {
	// Should not panic or alter unrelated strings.
	r := redact.New("")
	input := "clean string"
	if got := r.Sanitise(input); got != input {
		t.Errorf("expected %q unchanged, got %q", input, got)
	}
}

func TestSanitiseSlice(t *testing.T) {
	r := redact.New("tok")
	ss := []string{"value-tok-end", "normal"}
	out := r.SanitiseSlice(ss)
	if out[0] != "value-***-end" {
		t.Errorf("expected redacted, got %q", out[0])
	}
	if out[1] != "normal" {
		t.Errorf("expected unchanged, got %q", out[1])
	}
}

func TestContainsSecret_True(t *testing.T) {
	r := redact.New("hunter2")
	if !r.ContainsSecret("password is hunter2 here") {
		t.Error("expected ContainsSecret to return true")
	}
}

func TestContainsSecret_False(t *testing.T) {
	r := redact.New("hunter2")
	if r.ContainsSecret("nothing sensitive") {
		t.Error("expected ContainsSecret to return false")
	}
}
