package docker

import (
	"testing"
)

func TestTrimLeadingSlash(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/mycontainer", "mycontainer"},
		{"mycontainer", "mycontainer"},
		{"/", ""},
		{"", ""},
		{"//double", "/double"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := trimLeadingSlash(tt.input)
			if got != tt.want {
				t.Errorf("trimLeadingSlash(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainerInfoFields(t *testing.T) {
	info := ContainerInfo{
		Name:  "api",
		Image: "nginx:1.25",
		State: "running",
	}

	if info.Name != "api" {
		t.Errorf("expected Name=api, got %s", info.Name)
	}
	if info.Image != "nginx:1.25" {
		t.Errorf("expected Image=nginx:1.25, got %s", info.Image)
	}
	if info.State != "running" {
		t.Errorf("expected State=running, got %s", info.State)
	}
}
