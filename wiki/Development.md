# Development Guide

This guide explains how to set up a development environment for contributing to phev2mqtt.

## Prerequisites

- **Go 1.16+** installed locally (for local development)
- **Docker and Docker Compose** (for containerized development)
- **VS Code** with Go extension (recommended)
- **Git** for version control
- **MQTT broker** for testing (Mosquitto recommended)

## Quick Start

### Clone the Repository

```bash
git clone https://github.com/stefanh12/phev2mqtt.git
cd phev2mqtt
```

### Configure Environment

```bash
cp .env.example .env
nano .env  # Edit with your MQTT broker details
```

### Choose Development Method

- **[Docker Development](#docker-development)** - Containerized environment
- **[Local Development](#local-development)** - Native Go development with debugging

---

## Docker Development

### Setup

1. **Ensure Docker is running**
2. **Configure `.env` file** with your MQTT broker settings
3. **Build and run:**

```bash
# Using docker-compose
docker-compose -f docker-compose.dev.yml up --build

# Or using VS Code task: "Docker: Rebuild and Run"
```

### VS Code Tasks for Docker

Available tasks (press `Ctrl+Shift+P` → "Tasks: Run Task"):

- **Docker: Build Development Image** - Build the Docker image
- **Docker: Run Development Container** - Start the container
- **Docker: Stop Development Container** - Stop the container
- **Docker: Rebuild and Run** - Rebuild and start in one step
- **Docker: View Logs** - Follow container logs

### Docker Development Workflow

1. **Make changes** to source code
2. **Rebuild and run:**
   ```bash
   docker-compose -f docker-compose.dev.yml up --build
   ```
3. **View logs:**
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f
   ```
4. **Stop:**
   ```bash
   docker-compose -f docker-compose.dev.yml down
   ```

### Docker Configuration

The development container:
- Mounts the workspace as a volume
- Runs with `NET_ADMIN` capability for routing
- Uses `docker-compose.dev.yml` configuration
- Environment variables from `.env` file

---

## Local Development

### Setup

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Build the project:**
   ```bash
   go build -o phev2mqtt .
   # Or use VS Code task: "Go: Build Local" (Ctrl+Shift+B)
   ```

3. **Run the MQTT client:**
   ```bash
   ./phev2mqtt client mqtt \
     --mqtt_server tcp://192.168.1.2:1883 \
     --mqtt_username phevmqttuser \
     --mqtt_password mqttuserpassword
   
   # Or use VS Code task: "Go: Run MQTT Client"
   ```

### VS Code Debugging

VS Code is configured with launch configurations for debugging:

#### Launch MQTT Client

Debug the MQTT client with full debugging support:

1. **Set breakpoints** in your code
2. **Press F5** or go to Run and Debug panel
3. **Select** "Launch MQTT Client"
4. **Start debugging**

Configuration in `.vscode/launch.json`:
```json
{
  "name": "Launch MQTT Client",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}",
  "args": [
    "client", "mqtt",
    "--mqtt_server", "tcp://192.168.1.2:1883",
    "--mqtt_username", "phevmqttuser",
    "--mqtt_password", "mqttuserpassword"
  ]
}
```

#### Launch Client Register

Debug the registration process:

1. **Select** "Launch Client Register" from debug configurations
2. **Start debugging**

This runs in registration mode for first-time vehicle setup.

### VS Code Tasks for Local Development

- **Go: Build Local** (default build task - `Ctrl+Shift+B` or `Cmd+Shift+B`)
- **Go: Run MQTT Client** - Build and run

---

## Project Structure

```
phev2mqtt/
├── cmd/                    # Command-line interface
│   ├── root.go            # Root command
│   ├── client.go          # Client command
│   ├── mqtt.go            # MQTT client implementation
│   ├── decode.go          # Decode commands
│   ├── emulator.go        # PHEV emulator
│   ├── file.go            # File operations
│   ├── hex.go             # Hex utilities
│   ├── pcap.go            # Packet capture
│   ├── register.go        # Registration
│   ├── reload.go          # Config reload
│   ├── set.go             # Set commands
│   └── watch.go           # File watcher
├── client/                 # PHEV client library
│   └── client.go
├── emulator/               # PHEV emulator
│   ├── car.go             # Car state
│   ├── connection.go      # Network handling
│   └── state_machine.go   # State machine
├── protocol/               # PHEV protocol implementation
│   ├── message.go         # Message handling
│   ├── message_test.go    # Tests
│   ├── raw.go             # Raw protocol
│   ├── raw_test.go        # Tests
│   ├── settings.go        # Settings
│   └── README.md          # Protocol docs
├── unraid/                 # Unraid deployment
│   ├── Dockerfile
│   ├── entrypoint.sh
│   ├── phev2mqtt.xml      # Unraid template
│   └── .env.example
├── docker-compose.yml      # Production compose
├── docker-compose.dev.yml  # Development compose
├── Dockerfile.dev          # Development Dockerfile
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── main.go                 # Application entry point
└── .env.example            # Example configuration
```

---

## Making Changes

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting: `go fmt ./...`
- Run linters: `golint ./...` (if installed)
- Write tests for new features
- Document exported functions

### Testing

Run tests:
```bash
# All tests
go test ./...

# Specific package
go test ./protocol

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Building for Production

```bash
# Build for current platform
go build -o phev2mqtt .

# Build for Linux (for Docker)
GOOS=linux GOARCH=amd64 go build -o phev2mqtt .

# Build with optimizations
go build -ldflags="-s -w" -o phev2mqtt .
```

---

## Testing with PHEV

### Using the Emulator

phev2mqtt includes a PHEV emulator for testing without a real vehicle:

```bash
# Start emulator
./phev2mqtt emulator

# In another terminal, connect client
./phev2mqtt client mqtt --mqtt_server tcp://localhost:1883 ...
```

The emulator simulates PHEV responses for testing.

### Testing with Real Vehicle

1. **Ensure WiFi connection** to your PHEV
2. **Set registration mode** if first time:
   ```bash
   phev_register=true
   ```
3. **Run with debug logging:**
   ```bash
   log_level=debug
   ```
4. **Monitor logs** for connection status and errors

### Network Configuration

For testing, you may need to add a route to the PHEV network:

```bash
# Linux/macOS
sudo ip route add 192.168.8.0/24 via 192.168.1.1

# Or set in .env
route_add=192.168.1.1
```

---

## Debugging

### Enable Debug Logging

```bash
# In .env file
log_level=debug
```

### Common Debug Scenarios

**Connection Issues:**
```bash
# Check MQTT connection
docker logs phev2mqtt | grep "MQTT"

# Check PHEV connection
docker logs phev2mqtt | grep "PHEV"

# Check routing
docker exec phev2mqtt ip route
```

**Protocol Debugging:**
```bash
# Decode hex packets
./phev2mqtt decode <hex_data>

# Watch raw packets
./phev2mqtt pcap <pcap_file>
```

**Configuration Reload:**
```bash
# Watch for config changes
docker logs phev2mqtt | grep "reload"
```

---

## Contributing

### Workflow

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Write tests** for new functionality
5. **Run tests:** `go test ./...`
6. **Commit with clear messages:**
   ```bash
   git commit -m "Add feature: description"
   ```
7. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```
8. **Create a Pull Request** on GitHub

### Pull Request Guidelines

- **Clear description** of changes and motivation
- **Tests included** for new features
- **Documentation updated** if needed
- **Code follows** Go conventions
- **No breaking changes** without discussion
- **Passes all tests**

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks

**Example:**
```
feat: Add hot-reload configuration support

- Implement file watcher for .env changes
- Add validation for reloadable settings
- Update documentation

Closes #123
```

---

## Building Container Images

### Development Image

```bash
docker build -f Dockerfile.dev -t phev2mqtt:dev .
```

### Production Image

```bash
docker build -t phev2mqtt:latest .
```

### Multi-platform Build

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t ghcr.io/stefanh12/phev2mqtt:latest .
```

---

## Tools and Utilities

### Included Commands

**Decode hex data:**
```bash
./phev2mqtt decode F6010000F7
```

**Watch packets:**
```bash
./phev2mqtt pcap capture.pcap
```

**File operations:**
```bash
./phev2mqtt file <command>
```

**Register with PHEV:**
```bash
./phev2mqtt register
```

**Set values:**
```bash
./phev2mqtt set <register> <value>
```

---

## Troubleshooting Development

### Build Errors

**Missing dependencies:**
```bash
go mod download
go mod tidy
```

**Outdated dependencies:**
```bash
go get -u ./...
go mod tidy
```

### Docker Issues

**Permission denied:**
```bash
# On Linux, add user to docker group
sudo usermod -aG docker $USER
# Then log out and back in
```

**Port conflicts:**
```bash
# Check what's using the port
lsof -i :8080
# Kill or stop conflicting service
```

### VS Code Issues

**Go extension not working:**
- Install Go tools: Command Palette → "Go: Install/Update Tools"
- Check Go path: `which go`
- Restart VS Code

**Debugger not attaching:**
- Verify launch configuration
- Check Go version compatibility
- Ensure delve is installed: `go install github.com/go-delve/delve/cmd/dlv@latest`

---

## Resources

### Documentation

- [Protocol Documentation](Protocol-Documentation) - PHEV protocol details
- [Configuration](Configuration) - Configuration reference
- [SECURITY_AUDIT.md](../SECURITY_AUDIT.md) - Security considerations

### External Resources

- [buxtronix/phev2mqtt](https://github.com/buxtronix/phev2mqtt) - Original project
- [CodeCutterUK/phev2mqtt](https://github.com/CodeCutterUK/phev2mqtt) - Enhanced version
- [phev-remote](https://github.com/phev-remote) - Protocol research
- [Go Documentation](https://golang.org/doc/) - Official Go docs
- [MQTT Specification](https://mqtt.org/mqtt-specification/) - MQTT protocol

---

## Getting Help

- **GitHub Issues** - Report bugs or request features
- **GitHub Discussions** - Ask questions and share ideas
- **Pull Requests** - Contribute improvements

## Next Steps

- [Protocol Documentation](Protocol-Documentation) - Understand the PHEV protocol
- [Security Best Practices](Security-Best-Practices) - Security considerations for development
- [Troubleshooting](Troubleshooting) - Common issues and solutions
