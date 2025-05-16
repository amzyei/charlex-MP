# CharlexMP

A simple Unix-like management panel with a web UI to run shell commands and manage the OS.

## Features

- Run shell commands via web interface
- View command output and errors
- Change directories
- Lightweight Go implementation

## Installation & Usage

1. Install Go (v1.16+).
2. Clone repo and build:
   ```bash
   git clone https://github.com/amzyei/CharlexMP.git
   cd CharlexMP
   go build -o charlexMP main.go
   ```
3. Run the server:
   ```bash
   ./charlexMP
   ```
4. Open browser at `http://localhost:8080/` and use `/execute` to run commands.

## Security

For local or trusted use only. No authentication; do not expose publicly.

