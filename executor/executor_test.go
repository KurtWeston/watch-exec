package executor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		executable string
		args       []string
		wantExe    string
		wantArgs   int
	}{
		{
			name:       "valid executable",
			executable: "echo",
			args:       []string{"hello", "world"},
			wantArgs:   2,
		},
		{
			name:       "empty executable",
			executable: "",
			args:       []string{"arg1"},
			wantExe:    "",
			wantArgs:   1,
		},
		{
			name:       "whitespace in args",
			executable: "echo",
			args:       []string{"  ", "valid", ""},
			wantArgs:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(tt.executable, tt.args)
			if exec == nil {
				t.Fatal("New() returned nil")
			}
			if tt.wantExe != "" && exec.executable != tt.wantExe {
				t.Errorf("executable = %v, want %v", exec.executable, tt.wantExe)
			}
			if len(exec.args) != tt.wantArgs {
				t.Errorf("args length = %v, want %v", len(exec.args), tt.wantArgs)
			}
		})
	}
}

func TestValidateExecutable(t *testing.T) {
	tests := []struct {
		name string
		exe  string
		want bool
	}{
		{"valid command", "echo", true},
		{"empty string", "", false},
		{"whitespace only", "  ", false},
		{"nonexistent", "nonexistent_cmd_12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateExecutable(tt.exe)
			if tt.want && result == "" {
				t.Errorf("validateExecutable() = empty, want non-empty")
			}
			if !tt.want && result != "" {
				t.Errorf("validateExecutable() = %v, want empty", result)
			}
		})
	}
}

func TestSanitizeArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{"normal args", []string{"a", "b", "c"}, 3},
		{"with whitespace", []string{"  a  ", "b", "  "}, 2},
		{"empty slice", []string{}, 0},
		{"all empty", []string{"", "  ", ""}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeArgs(tt.args)
			if len(got) != tt.want {
				t.Errorf("sanitizeArgs() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		executable string
		args       []string
		wantErr    bool
	}{
		{
			name:       "successful command",
			executable: "echo",
			args:       []string{"test"},
			wantErr:    false,
		},
		{
			name:       "empty executable",
			executable: "",
			args:       []string{},
			wantErr:    true,
		},
		{
			name:       "nonexistent command",
			executable: "nonexistent_cmd_xyz",
			args:       []string{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(tt.executable, tt.args)
			err := exec.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
