# Quick Start Guide

Get phev2mqtt up and running in minutes!

## Prerequisites Checklist

Before you begin, ensure you have:

- [ ] Mitsubishi Outlander PHEV (MY18 or newer recommended)
- [ ] MQTT broker running (e.g., Mosquitto)
- [ ] Docker and Docker Compose installed
- [ ] PHEV WiFi credentials (SSID and password)
- [ ] MQTT broker credentials

## 5-Minute Setup

### Step 1: Download Configuration Files (1 minute)

```bash
# Create directory
mkdir -p ~/phev2mqtt && cd ~/phev2mqtt

# Download files
wget https://raw.githubusercontent.com/stefanh12/phev2mqtt/main/docker-compose.yml
wget https://raw.githubusercontent.com/stefanh12/phev2mqtt/main/unraid/.env.example
```

### Step 2: Configure (2 minutes)

```bash
# Copy example to .env
cp .env.example .env

# Edit with your details
nano .env
```

**Minimal required config:**
```bash
# MQTT Settings
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_strong_password_here

# Network 
route_add=192.168.1.1

# Operation
phev_register=false
log_level=info
```

**Save and exit** (`Ctrl+X`, then `Y`, then `Enter`)

### Step 3: Start Container (1 minute)

```bash
docker-compose up -d
```

### Step 4: Verify (1 minute)

```bash
# Check logs
docker-compose logs -f

# Look for:
# ✅ "Connected to MQTT broker"
# ✅ "Waiting for PHEV connection" or "Connected to PHEV"
```

**Press `Ctrl+C` to exit logs**

### Step 5: Connect to PHEV WiFi

Ensure your Docker host (or WiFi bridge) is connected to your PHEV's WiFi network.

**Done!** Check Home Assistant - your PHEV device should appear automatically!

---

## First-Time Setup (Registration Required)

If this is your **first time** connecting to your PHEV, you need to register:

### 1. Put PHEV in Registration Mode

Consult your vehicle manual for the exact procedure (typically involves the remote app screen).

### 2. Enable Registration Mode

```bash
# Edit .env file
nano .env

# Set this line:
phev_register=true

# Save and restart
docker-compose restart
```

### 3. Follow Prompts

Watch the logs for registration prompts:
```bash
docker-compose logs -f
```

Follow the on-screen instructions and enter the security code from your PHEV dashboard.

### 4. Switch Back to Normal Mode

After successful registration:

```bash
# Edit .env file
nano .env

# Set this line:
phev_register=false

# Save and restart
docker-compose restart
```

---

## Unraid Quick Start

### Method 1: Community Applications

1. Open Unraid Web UI
2. Go to **Apps** tab
3. Search for "phev2mqtt"
4. Click **Install**
5. Configure settings in template
6. Click **Apply**

### Method 2: Manual Setup

1. Navigate to `/mnt/user/appdata/phev2mqtt/`
2. Create `.env` file with your configuration
3. Install from template or Docker tab
4. Start container

See [Installation](Installation#unraid-installation) for detailed instructions.

---

## Home Assistant Integration

**Auto-discovery is enabled by default!**

Just wait a few seconds after phev2mqtt connects to your PHEV:

1. Go to **Settings** → **Devices & Services** → **MQTT**
2. Look for "Mitsubishi Outlander PHEV" in discovered devices
3. Click to see all entities

**Optional:** Set your VIN for immediate discovery:
```bash
# In .env file
vehicle_vin=JA4J24A58KZ123456
```

See [Home Assistant Integration](Home-Assistant-Integration) for dashboards and automations.

---

## Common Issues

### "Can't connect to MQTT broker"

**Check:**
- MQTT broker IP and port correct
- Credentials valid
- Firewall allows connection

**Test:**
```bash
mosquitto_pub -h 192.168.1.2 -u phevmqttuser -P password -t "test" -m "hello"
```

### "Waiting for PHEV connection" (never connects)

**Check:**
- WiFi connected to PHEV network
- Can ping 192.168.8.46
- Route to PHEV network configured (`route_add` in .env)
- PHEV WiFi is active (vehicle in range)

**Test:**
```bash
ping 192.168.8.46
```

### "Not registered" or "Authentication failed"

**Solution:**
- Run first-time registration (see above)
- Set `phev_register=true`, restart, follow prompts
- After registration, set `phev_register=false`

### Entities not appearing in Home Assistant

**Check:**
- MQTT integration configured in Home Assistant
- phev2mqtt connected to PHEV (check logs)
- Wait a few minutes for discovery

**Force refresh:**
```bash
docker-compose restart
```

See [Troubleshooting](Troubleshooting) for more solutions.

---

## Next Steps

Now that you're up and running:

### Essential Reading

- **[Configuration](Configuration)** - Customize behavior, enable WiFi management
- **[Home Assistant Integration](Home-Assistant-Integration)** - Dashboards and automations
- **[Security Best Practices](Security-Best-Practices)** - Secure your installation

### Advanced Setup

- **[MikroTik Integration](MikroTik-Integration)** - Set up dedicated WiFi bridge
- **[WiFi Management](WiFi-Management)** - Power saving and reliability
- **[Protocol Documentation](Protocol-Documentation)** - Understand PHEV communication

### Need Help?

- **[Troubleshooting](Troubleshooting)** - Common issues and solutions
- **[GitHub Issues](https://github.com/stefanh12/phev2mqtt/issues)** - Report bugs
- **[GitHub Discussions](https://github.com/stefanh12/phev2mqtt/discussions)** - Ask questions

---

## Configuration Examples

### Basic Setup

```bash
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!
route_add=192.168.1.1
log_level=info
update_interval=5m
```

### With MikroTik WiFi Management

```bash
# Basic config...
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!

# Remote WiFi restart
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart

# Power save mode
remote_wifi_power_save_enabled=true
remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi
remote_wifi_enable_message={"wifi": "enable"}
remote_wifi_disable_message={"wifi": "disable"}
update_interval=10m
```

### Debug Mode

```bash
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!
log_level=debug
update_interval=15m
```

---

## Useful Commands

### Container Management

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Restart
docker-compose restart

# View logs
docker-compose logs -f

# Update to latest version
docker-compose pull && docker-compose up -d
```

### MQTT Testing

```bash
# Publish test message
mosquitto_pub -h 192.168.1.2 -u user -P pass -t "test" -m "hello"

# Subscribe to phev2mqtt messages
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "phev/#" -v

# Subscribe to Home Assistant discovery
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "homeassistant/#" -v
```

### Debugging

```bash
# Check container status
docker ps | grep phev2mqtt

# Check routes in container
docker exec phev2mqtt ip route

# Test PHEV connection from container
docker exec phev2mqtt ping -c 3 192.168.8.46

# Access container shell
docker exec -it phev2mqtt /bin/sh
```

---

## Updating

Keep phev2mqtt up to date for bug fixes and new features:

```bash
# Pull latest image
docker-compose pull

# Restart with new image
docker-compose up -d

# Check version in logs
docker-compose logs | head -20
```

**Subscribe to releases:**
- [GitHub Releases](https://github.com/stefanh12/phev2mqtt/releases)
- Watch repository for notifications

---

## Support the Project

If you find phev2mqtt useful:

- ⭐ Star the [GitHub repository](https://github.com/stefanh12/phev2mqtt)
- 🐛 Report bugs and issues
- 💡 Suggest features
- 🔧 Contribute code improvements
- 📚 Improve documentation
- 💬 Help others in Discussions

---

**Happy PHEVing! 🚗⚡**
