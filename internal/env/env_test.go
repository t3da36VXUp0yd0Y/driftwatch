package env

import (
	"testing"
)

func TestCompare_NoDiff(t *testing.T) {
	declared := map[string]string{"PORT": "8080", "ENV": "prod"}
	actual := map[string]string{"PORT": "8080", "ENV": "prod", "EXTRA": "ignored"}
	diffs := Compare(declared, actual)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompare_ValueMismatch(t *testing.T) {
	declared := map[string]string{"PORT": "8080"}
	actual := map[string]string{"PORT": "9090"}
	diffs := Compare(declared, actual)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "PORT" || diffs[0].Expected != "8080" || diffs[0].Actual != "9090" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompare_MissingKey(t *testing.T) {
	declared := map[string]string{"SECRET": "abc"}
	actual := map[string]string{}
	diffs := Compare(declared, actual)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if !diffs[0].Missing {
		t.Error("expected Missing=true")
	}
}

func TestCompare_MultipleDiffs(t *testing.T) {
	declared := map[string]string{"A": "1", "B": "2", "C": "3"}
	actual := map[string]string{"A": "1", "B": "99"}
	diffs := Compare(declared, actual)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
}

func TestDiff_String_Missing(t *testing.T) {
	d := Diff{Key: "FOO", Expected: "bar", Missing: true}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestDiff_String_Mismatch(t *testing.T) {
	d := Diff{Key: "FOO", Expected: "bar", Actual: "baz"}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestParseEnvSlice(t *testing.T) {
	slice := []string{"KEY=value", "PORT=8080", "INVALID"}
	m := ParseEnvSlice(slice)
	if m["KEY"] != "value" {
		t.Errorf("expected value, got %q", m["KEY"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", m["PORT"])
	}
	if _, ok := m["INVALID"]; ok {
		t.Error("INVALID should not be in map")
	}
}

func TestParseEnvSlice_Empty(t *testing.T) {
	m := ParseEnvSlice(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}
