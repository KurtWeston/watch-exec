package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/watch-exec/executor"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Patterns: []string{tmpDir},
				Debounce: 100 * time.Millisecond,
				MaxDepth: 5,
			},
			wantErr: false,
		},
		{
			name: "nonexistent path",
			config: Config{
				Patterns: []string{"/nonexistent/path/xyz"},
				Debounce: 100 * time.Millisecond,
				MaxDepth: 5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if w != nil {
				w.Close()
			}
		})
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		name    string
		ignore  []string
		path    string
		wantIgn bool
	}{
		{"dot file", []string{}, ".git", true},
		{"normal file", []string{}, "file.txt", false},
		{"pattern match", []string{"*.log"}, "test.log", true},
		{"pattern no match", []string{"*.log"}, "test.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Watcher{
				config: Config{Ignore: tt.ignore},
			}
			got := w.shouldIgnore(tt.path)
			if got != tt.wantIgn {
				t.Errorf("shouldIgnore() = %v, want %v", got, tt.wantIgn)
			}
		})
	}
}

func TestAddPath(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)

	tests := []struct {
		name      string
		path      string
		recursive bool
		maxDepth  int
		wantErr   bool
	}{
		{"add directory", tmpDir, false, 10, false},
		{"add recursive", tmpDir, true, 10, false},
		{"max depth exceeded", tmpDir, true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := New(Config{
				Patterns:  []string{},
				Recursive: tt.recursive,
				MaxDepth:  tt.maxDepth,
			})
			defer w.Close()

			err := w.addPath(tt.path, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("addPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWatch(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	w, err := New(Config{
		Patterns: []string{tmpDir},
		Debounce: 50 * time.Millisecond,
		MaxDepth: 5,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	exec := executor.New("echo", []string{"test"})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		os.WriteFile(testFile, []byte("test"), 0644)
	}()

	err = w.Watch(ctx, exec)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Watch() error = %v", err)
	}
}

func TestClose(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := New(Config{
		Patterns: []string{tmpDir},
		Debounce: 100 * time.Millisecond,
		MaxDepth: 5,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
