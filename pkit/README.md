# pkit

A minimal personal CLI toolkit written in Go.

## Installation

```bash
go build -o pkit
```

## Usage

```bash
# Show help
./pkit --help

# Show version
./pkit --version

# Delete spam messages
./pkit delete-spam

# Download Google Drive folder
./pkit download-drive "https://drive.google.com/drive/folders/FOLDER_ID"

# Use custom config
./pkit --config custom-config.yaml delete-spam
```

## Configuration

Configuration is managed via `config.yaml`. Copy `.env.sample` to `.env` for environment variables.

## Build

```bash
# Build for current platform
make build

# Build for all platforms
make all

# Bump version
make version-bump

# Clean build artifacts
make clean
```