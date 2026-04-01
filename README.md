# watch-exec

Watch files for changes and execute commands with intelligent debouncing and pattern filtering

## Features

- Recursive directory watching with configurable depth
- Glob pattern filtering for include/exclude paths
- Intelligent debouncing to group rapid file changes
- Execute arbitrary shell commands on file changes
- Colorized output showing watched paths and execution status
- Clear feedback on command success/failure with exit codes
- Initial execution on startup (optional flag)
- Graceful shutdown with cleanup on SIGINT/SIGTERM
- Low CPU usage with efficient event handling
- Support for multiple watch patterns in single invocation
- Configurable debounce delay (default 300ms)
- Verbose mode showing all detected file events

## How to Use

Use this project when you need to:

- Quickly solve problems related to watch-exec
- Integrate go functionality into your workflow
- Learn how go handles common patterns

## Installation

```bash
# Clone the repository
git clone https://github.com/KurtWeston/watch-exec.git
cd watch-exec

# Install dependencies
go build
```

## Usage

```bash
./main
```

## Built With

- go

## Dependencies

- `github.com/fsnotify/fsnotify`
- `github.com/spf13/cobra`
- `github.com/fatih/color`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
