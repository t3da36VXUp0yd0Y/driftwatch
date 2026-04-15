package version_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/version"
)

func TestGet_DefaultValues(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("expected non-empty Version")
	}
	if info.Commit == "" {
		t.Error("expected non-empty Commit")
	}
	if info.Date == "" {
		t.Error("expected non-empty Date")
	}
}

func TestGet_ReturnsInfo(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc1234"
	version.Date = "2024-01-15T10:00:00Z"

	info := version.Get()

	if info.Version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %s", info.Version)
	}
	if info.Commit != "abc1234" {
		t.Errorf("expected commit abc1234, got %s", info.Commit)
	}
	if info.Date != "2024-01-15T10:00:00Z" {
		t.Errorf("expected date 2024-01-15T10:00:00Z, got %s", info.Date)
	}
}

func TestInfo_String_ContainsVersion(t *testing.T) {
	version.Version = "0.5.0"
	version.Commit = "deadbeef"
	version.Date = "2024-06-01T00:00:00Z"

	info := version.Get()
	s := info.String()

	if !strings.Contains(s, "0.5.0") {
		t.Errorf("expected string to contain version, got: %s", s)
	}
	if !strings.Contains(s, "deadbeef") {
		t.Errorf("expected string to contain commit, got: %s", s)
	}
	if !strings.Contains(s, "driftwatch") {
		t.Errorf("expected string to contain 'driftwatch', got: %s", s)
	}
}

func TestInfo_String_Format(t *testing.T) {
	version.Version = "2.0.0"
	version.Commit = "cafebabe"
	version.Date = "2024-12-01T00:00:00Z"

	info := version.Get()
	s := info.String()

	expected := "driftwatch 2.0.0 (commit: cafebabe, built: 2024-12-01T00:00:00Z)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
