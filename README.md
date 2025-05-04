# Smol Docker

A minimal container runtime written in Go, designed for educational purposes and understanding container internals.

## Features

- Pull and run Docker images
- Cross-platform support (Linux and macOS)
- Simple command-line interface
- Basic container isolation on Linux using chroot and mount namespaces
- Image extraction and management

## Prerequisites

- Go 1.23.2 or later
- Docker (for pulling images)
- Linux or macOS operating system

## Installation

1. Clone the repository:
```bash
git clone https://github.com/smol-go/smol-docker.git
cd smol-docker
```

2. Build the binary:
```bash
make build
```

## Usage

### Pull an Image
```bash
./smol-docker pull <image-name>
```

### Run a Container
```bash
./smol-docker run <image-name> [command]
```

If no command is specified, the default command from the image will be used.

## Platform Support

### Linux
- Full container isolation using chroot and mount namespaces
- Proper process isolation
- Filesystem isolation

### macOS
- Limited container support
- Runs in a simulated environment
- Note: Cannot run Linux binaries directly on macOS. Use Docker Desktop or a Linux VM for full container support.

## Project Structure

- `main.go` - Core container runtime implementation
- `linux.go` - Linux-specific container isolation features
- `pull.sh` - Script to extract Docker images
- `Makefile` - Build and distribution management

## Limitations

- Basic container isolation (no network, cgroups, or advanced security features)
- Limited to running single processes
- No image layer management
- No container networking
- macOS support is limited to basic environment simulation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.