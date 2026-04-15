package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftwatch/internal/output"
)

func TestNew_Stdout(t *testing.T) {
	w, err := output.New(output.Options{Dest: output.Stdout})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write to stdout failed: %v", err)
	}
}

func TestNew_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w, err := output.New(output.Options{Dest: output.File, FilePath: path})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = w.Write([]byte("hello drift"))
	if err != nil {
		t.Fatalf("Write to file failed: %v", err)
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if string(data) != "hello drift" {
		t.Errorf("expected %q, got %q", "hello drift", string(data))
	}
}

func TestNew_File_EmptyPath(t *testing.T) {
	_, err := output.New(output.Options{Dest: output.File, FilePath: ""})
	if err == nil {
		t.Fatal("expected error for empty file path, got nil")
	}
}

func TestNew_UnknownDestination(t *testing.T) {
	_, err := output.New(output.Options{Dest: output.Destination(99)})
	if err == nil {
		t.Fatal("expected error for unknown destination, got nil")
	}
}
