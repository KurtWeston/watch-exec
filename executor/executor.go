package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Executor struct {
	executable string
	args       []string
}

func New(executable string, args []string) *Executor {
	return &Executor{
		executable: validateExecutable(executable),
		args:       sanitizeArgs(args),
	}
}

func validateExecutable(exe string) string {
	exe = strings.TrimSpace(exe)
	if exe == "" {
		return ""
	}

	if filepath.IsAbs(exe) {
		if _, err := os.Stat(exe); err == nil {
			return exe
		}
	}

	if path, err := exec.LookPath(exe); err == nil {
		return path
	}

	return exe
}

func sanitizeArgs(args []string) []string {
	sanitized := make([]string, 0, len(args))
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg != "" {
			sanitized = append(sanitized, arg)
		}
	}
	return sanitized
}

func (e *Executor) Run() error {
	if e.executable == "" {
		return fmt.Errorf("no executable specified")
	}

	start := time.Now()
	color.Cyan("[%s] Executing: %s %s\n",
		start.Format("15:04:05"),
		e.executable,
		strings.Join(e.args, " "))

	cmd := exec.Command(e.executable, e.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			color.Red("[%s] Command failed with exit code %d (took %v)\n",
				time.Now().Format("15:04:05"),
				exitErr.ExitCode(),
				duration)
			return fmt.Errorf("exit code %d", exitErr.ExitCode())
		}
		color.Red("[%s] Command failed: %v (took %v)\n",
			time.Now().Format("15:04:05"),
			err,
			duration)
		return err
	}

	color.Green("[%s] Command succeeded (took %v)\n",
		time.Now().Format("15:04:05"),
		duration)
	return nil
}
