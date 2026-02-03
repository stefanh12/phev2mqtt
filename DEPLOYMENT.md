# Production Deployment with Docker Compose

This guide shows how to deploy phev2mqtt using docker-compose with environment variables stored in a `.env` file.

## Setup

1. **Create your `.env` file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your actual credentials:**
   ```bash
   MQTT_SERVER=192.168.1.2:1883
   MQTT_USER=phevmqttuser
   MQTT_PASSWORD=your_actual_password
   PHEV_REGISTER=false
   DEBUG=false
   ROUTE_ADD=192.168.1.1
   ```

3. **Deploy:**
   ```bash
   docker-compose up -d
   ```

4. **View logs:**
   ```bash
   docker-compose logs -f
   ```

5. **Stop:**
   ```bash
   docker-compose down
   ```

## For Unraid Users

Unraid supports two methods for configuration:

### Method 1: Using .env File (Recommended for Security)

1. Install from Community Applications or import the template from `unraid/phev2mqtt.xml`
2. The template will map `/mnt/user/appdata/phev2mqtt` to `/config` in the container
3. Create your `.env` file:
   - Navigate to `/mnt/user/appdata/phev2mqtt/` on your Unraid server
   - Copy `unraid/.env.example` to `.env`
   - Edit `.env` with your credentials (use nano, vi, or Unraid's built-in editor)
4. Start the container - it will automatically load settings from `/config/.env`

**Benefits:**
- Passwords stored in a file, not in container configuration
- Easy to backup and manage
- Can be edited without recreating the container

### Method 2: Environment Variables via Unraid Web UI

Configure variables directly in the Unraid Docker template:
   - mqtt_server
   - mqtt_user
   - mqtt_password (will be masked in UI)
   - phev_register
   - debug
   - route_add

**Note:** If both methods are used, the `.env` file takes precedence.

### Alternative: Manual Docker on Unraid

If you want to use `.env` on Unraid without the template:

1. SSH into your Unraid server
2. Create a directory: `mkdir -p /mnt/user/appdata/phev2mqtt`
3. Copy `.env.example` to `/mnt/user/appdata/phev2mqtt/.env`
4. Edit the `.env` file with your credentials
5. Copy `docker-compose.yml` to the same directory
6. Run: `cd /mnt/user/appdata/phev2mqtt && docker-compose up -d`

**Note:** This method bypasses Unraid's docker management and won't appear in the Unraid Docker tab.

## Environment Variables

All environment variables with defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `MQTT_SERVER` | `192.168.1.2:1883` | MQTT broker address and port |
| `MQTT_USER` | `phevmqttuser` | MQTT username |
| `MQTT_PASSWORD` | *(required)* | MQTT password |
| `PHEV_REGISTER` | `false` | Set to `true` for registration mode |
| `DEBUG` | `false` | Enable debug mode |
| `ROUTE_ADD` | `192.168.1.1` | Gateway IP for routing to PHEV |

## Updating

```bash
docker-compose pull
docker-compose up -d
```

## Security Notes

- Never commit your `.env` file to git (it's in `.gitignore`)
- Use strong, unique passwords for MQTT
- Consider using MQTT over TLS in production
- Restrict MQTT user permissions to only what's needed
