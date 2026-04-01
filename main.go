package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user/watch-exec/executor"
	"github.com/user/watch-exec/watcher"
)

var (
	patterns      []string
	ignorePattern []string
	debounceMs    int
	recursive     bool
	initialRun    bool
	verbose       bool
	maxDepth      int
)

var rootCmd = &cobra.Command{
	Use:   "watch-exec [flags] -- <command> [args...]",
	Short: "Watch files and execute commands on changes",
	Long:  "Watch files for changes and execute commands with intelligent debouncing and pattern filtering.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	rootCmd.Flags().StringSliceVarP(&patterns, "pattern", "p", []string{"*"}, "Glob patterns to watch")
	rootCmd.Flags().StringSliceVarP(&ignorePattern, "ignore", "i", []string{}, "Patterns to ignore")
	rootCmd.Flags().IntVarP(&debounceMs, "debounce", "d", 300, "Debounce delay in milliseconds")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Watch directories recursively")
	rootCmd.Flags().BoolVar(&initialRun, "initial", false, "Run command on startup")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show all file events")
	rootCmd.Flags().IntVar(&maxDepth, "max-depth", 10, "Maximum recursion depth")
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command required after --")
	}

	exec := executor.New(args[0], args[1:])

	config := watcher.Config{
		Patterns:   patterns,
		Ignore:     ignorePattern,
		Debounce:   time.Duration(debounceMs) * time.Millisecond,
		Recursive:  recursive,
		Verbose:    verbose,
		MaxDepth:   maxDepth,
	}

	w, err := watcher.New(config)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.Close()

	color.Green("Watching: %s\n", strings.Join(patterns, ", "))
	color.Cyan("Command: %s %s\n", args[0], strings.Join(args[1:], " "))

	if initialRun {
		color.Yellow("Running initial command...\n")
		if err := exec.Run(); err != nil {
			color.Red("Initial run failed: %v\n", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		color.Yellow("\nShutting down...\n")
		cancel()
	}()

	return w.Watch(ctx, exec)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
