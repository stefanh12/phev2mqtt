# Development Setup

This guide explains how to develop and test phev2mqtt from VS Code.

## Prerequisites

- Docker and Docker Compose installed
- VS Code with Go extension installed
- Go 1.16+ installed locally (for local development without Docker)

## Quick Start

### Option 1: Docker Development

1. **Configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your MQTT broker details
   ```

2. **Build and run in VS Code:**
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
   - Type "Tasks: Run Task"
   - Select "Docker: Rebuild and Run"

3. **View logs:**
   - Run task: "Docker: View Logs"

4. **Stop container:**
   - Run task: "Docker: Stop Development Container"

### Option 2: Local Development

1. **Build locally:**
   - Press `Ctrl+Shift+B` (or `Cmd+Shift+B` on Mac) to build
   - Or run task: "Go: Build Local"

2. **Run with debugger:**
   - Press `F5` or go to Run and Debug panel
   - Select "Launch MQTT Client"
   - Set breakpoints in your code
   - Start debugging

3. **Run MQTT client from task:**
   - Run task: "Go: Run MQTT Client"

## Available VS Code Tasks

- **Docker: Build Development Image** - Build the Docker image
- **Docker: Run Development Container** - Start the container
- **Docker: Stop Development Container** - Stop the container
- **Docker: Rebuild and Run** - Rebuild and start container
- **Docker: View Logs** - Follow container logs
- **Go: Build Local** - Build binary locally
- **Go: Run MQTT Client** - Build and run MQTT client locally

## Debug Configurations

### Launch MQTT Client
Debug the MQTT client directly from VS Code with full debugging support.

### Launch Client Register
Debug the client registration process.

### Docker: Attach to Container
Attach debugger to running container (requires delve in container).

## Environment Variables

Edit these in `docker-compose.dev.yml` or your `.env` file:

- `mqtt_server` - MQTT broker address (e.g., 192.168.1.2:1883)
- `mqtt_user` - MQTT username
- `mqtt_password` - MQTT password
- `phev_register` - Set to "true" to run registration mode
- `debug` - Set to "true" to enable debug mode
- `route_add` - Gateway IP for routing to PHEV network

## Modifying Code

1. Make changes to the code
2. If using Docker:
   - The container volume mounts the workspace
   - Rebuild with: "Docker: Rebuild and Run" task
3. If using local Go:
   - Press `F5` to rebuild and debug
   - Or run "Go: Build Local" task

## Network Configuration

The container runs with `NET_ADMIN` capability to allow adding routes. If you need to route to the PHEV network (192.168.8.0/24), set the `route_add` variable to your gateway IP.

## Troubleshooting

- **Connection issues**: Check MQTT broker address and credentials
- **Route issues**: Verify `route_add` is set correctly
- **Build errors**: Run `go mod download` to ensure dependencies are installed
- **Docker issues**: Check `docker-compose -f docker-compose.dev.yml logs`

## Testing Changes

After making changes:

1. Build: `Ctrl+Shift+B` (or `Cmd+Shift+B`)
2. Run tests: `go test ./...`
3. Run locally or in Docker to verify functionality
