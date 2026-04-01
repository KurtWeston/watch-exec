package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/user/watch-exec/executor"
)

type Config struct {
	Patterns  []string
	Ignore    []string
	Debounce  time.Duration
	Recursive bool
	Verbose   bool
	MaxDepth  int
}

type Watcher struct {
	watcher  *fsnotify.Watcher
	config   Config
	debounce *time.Timer
	pending  bool
}

func New(config Config) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher: fsw,
		config:  config,
	}

	for _, pattern := range config.Patterns {
		if err := w.addPath(pattern, 0); err != nil {
			fsw.Close()
			return nil, err
		}
	}

	return w, nil
}

func (w *Watcher) addPath(path string, depth int) error {
	if depth > w.config.MaxDepth {
		return nil
	}

	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return nil
	}

	if !info.IsDir() {
		return w.watcher.Add(filepath.Dir(path))
	}

	if err := w.watcher.Add(path); err != nil {
		return err
	}

	if w.config.Recursive {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() && !w.shouldIgnore(entry.Name()) {
				subPath := filepath.Join(path, entry.Name())
				if err := w.addPath(subPath, depth+1); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (w *Watcher) shouldIgnore(name string) bool {
	for _, pattern := range w.config.Ignore {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return strings.HasPrefix(name, ".")
}

func (w *Watcher) Watch(ctx context.Context, exec *executor.Executor) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}
			if w.config.Verbose {
				color.Blue("Event: %s %s\n", event.Op, event.Name)
			}
			w.handleEvent(exec)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			color.Red("Watcher error: %v\n", err)
		}
	}
}

func (w *Watcher) handleEvent(exec *executor.Executor) {
	w.pending = true
	if w.debounce != nil {
		w.debounce.Stop()
	}
	w.debounce = time.AfterFunc(w.config.Debounce, func() {
		if w.pending {
			w.pending = false
			if err := exec.Run(); err != nil {
				color.Red("Execution failed: %v\n", err)
			}
		}
	})
}

func (w *Watcher) Close() error {
	if w.debounce != nil {
		w.debounce.Stop()
	}
	return w.watcher.Close()
}
