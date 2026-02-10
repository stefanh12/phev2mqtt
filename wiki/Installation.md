# Installation Guide

This guide covers installing phev2mqtt using Docker Compose or on Unraid.

## Prerequisites

Before you begin, ensure you have:

- ✅ Docker and Docker Compose installed (for Docker installation)
- ✅ An MQTT broker running (e.g., Mosquitto)
- ✅ Network connectivity to your Mitsubishi Outlander PHEV WiFi
- ✅ MQTT broker credentials (username and password)

## Installation Methods

Choose the installation method that best fits your setup:

1. **[Docker Compose](#docker-compose-installation)** - For general Docker environments
2. **[Unraid Template](#unraid-installation)** - For Unraid servers (recommended)
3. **[Unraid Manual](#unraid-manual-installation)** - Alternative Unraid setup

---

## Docker Compose Installation

### Step 1: Create Project Directory

```bash
mkdir -p ~/phev2mqtt
cd ~/phev2mqtt
```

### Step 2: Download Configuration Files

Download `docker-compose.yml` and `.env.example` from the repository:

```bash
# Download docker-compose.yml
wget https://raw.githubusercontent.com/stefanh12/phev2mqtt/main/docker-compose.yml

# Download .env.example
wget https://raw.githubusercontent.com/stefanh12/phev2mqtt/main/unraid/.env.example
```

### Step 3: Create Environment File

```bash
cp .env.example .env
nano .env  # or use your preferred editor
```

Edit the `.env` file with your actual credentials. **Required settings:**

```bash
# MQTT Configuration (REQUIRED)
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_strong_password_here

# Network Configuration
route_add=192.168.1.1

# PHEV Settings
phev_register=false
```

See [Configuration](Configuration) for all available options.

### Step 4: Start the Container

```bash
docker-compose up -d
```

### Step 5: Verify Installation

Check the logs to ensure successful connection:

```bash
docker-compose logs -f
```

You should see output similar to:
```
phev2mqtt starting...
Connected to MQTT broker at tcp://192.168.1.2:1883
Waiting for PHEV connection...
```

### Updating

To update to the latest version:

```bash
docker-compose pull
docker-compose up -d
```

---

## Unraid Installation

### Method 1: Community Applications (Recommended)

1. **Open Unraid Web UI** → Navigate to the **Apps** tab
2. **Search** for "phev2mqtt"
3. **Click Install** on the phev2mqtt template
4. **Configure** the template with your settings (see below)
5. **Apply** to start the container

### Method 2: Import Template

1. **Download** the template file from: `unraid/phev2mqtt.xml`
2. **Navigate to** Docker tab in Unraid
3. **Click** "Add Container"
4. **Click** "Template repositories" and add the GitHub URL
5. **Select** phev2mqtt from the list
6. **Configure** and **Apply**

### Unraid Configuration

The Unraid template creates `/mnt/user/appdata/phev2mqtt` for configuration storage.

#### Option A: Using .env File (Recommended)

This method keeps credentials secure and separate from the container configuration.

1. **Navigate** to `/mnt/user/appdata/phev2mqtt/` on your Unraid server
2. **Create** a `.env` file:
   ```bash
   cd /mnt/user/appdata/phev2mqtt
   nano .env
   ```
3. **Add your configuration:**
   ```bash
   mqtt_server=tcp://192.168.1.2:1883
   mqtt_username=phevmqttuser
   mqtt_password=your_strong_password
   route_add=192.168.1.1
   phev_register=false
   log_level=info
   ```
4. **Save** and restart the container

**Benefits:**
- Passwords stored securely in a file, not in container config
- Easy to backup and restore
- Can be edited without recreating the container
- Changes apply automatically via hot-reload

#### Option B: Environment Variables

Configure variables directly in the Unraid Docker template:

- `mqtt_server` - MQTT broker address
- `mqtt_username` - MQTT username
- `mqtt_password` - MQTT password (masked in UI)
- `route_add` - Gateway IP for routing
- `phev_register` - Set to `true` for registration mode
- `log_level` - Logging level (info, debug, etc.)

**Note:** If both methods are used, the `.env` file takes precedence.

### Unraid Manual Installation

For advanced users who want more control:

1. **SSH into your Unraid server**
2. **Create directory:**
   ```bash
   mkdir -p /mnt/user/appdata/phev2mqtt
   cd /mnt/user/appdata/phev2mqtt
   ```
3. **Create `.env` file** with your configuration
4. **Create `docker-compose.yml`:**
   ```yaml
   version: '3.8'
   services:
     phev2mqtt:
       image: ghcr.io/stefanh12/phev2mqtt:latest
       container_name: phev2mqtt
       restart: unless-stopped
       network_mode: host
       cap_add:
         - NET_ADMIN
       volumes:
         - /mnt/user/appdata/phev2mqtt/.env:/app/.env:ro
       environment:
         - TZ=America/New_York
   ```
5. **Start the container:**
   ```bash
   docker-compose up -d
   ```

**Note:** This method bypasses Unraid's Docker management and won't appear in the Docker tab.

---

## Post-Installation Steps

### 1. Verify MQTT Connection

Check that the container can connect to your MQTT broker:

```bash
# Docker Compose
docker-compose logs | grep "Connected to MQTT"

# Unraid
docker logs phev2mqtt | grep "Connected to MQTT"
```

### 2. Connect to PHEV WiFi

Ensure your Docker host (or WiFi bridge) is connected to your PHEV's WiFi network:

- **SSID**: REMOTE followed by your vehicle identifier
- **Default Gateway**: 192.168.8.46
- **Your IP**: Usually 192.168.8.47

### 3. Configure Home Assistant

If using Home Assistant, the devices should auto-discover. See [Home Assistant Integration](Home-Assistant-Integration) for details.

### 4. Set Up MikroTik Bridge (Optional)

For enhanced WiFi management, configure a MikroTik bridge. See [MikroTik Integration](MikroTik-Integration).

---

## First-Time Registration

If this is your first time connecting to your PHEV, you need to register:

1. **Set registration mode** in `.env`:
   ```bash
   phev_register=true
   ```
2. **Restart the container**
3. **Follow the on-screen prompts** in the logs
4. **Enter the security code** displayed on your PHEV dashboard
5. **After successful registration**, set back to:
   ```bash
   phev_register=false
   ```
6. **Restart the container**

See [Troubleshooting](Troubleshooting) if you encounter registration issues.

---

## Network Configuration

### Routing to PHEV Network

The PHEV uses the `192.168.8.0/24` network. If your Docker host needs routing:

```bash
# In .env file
route_add=192.168.1.1  # Your gateway IP
```

The container automatically adds the route:
```bash
ip route add 192.168.8.0/24 via <route_add>
```

### Using a WiFi Bridge

For separate VLAN or dedicated WiFi bridge (recommended):

1. **Configure your bridge** to connect to PHEV WiFi
2. **Set up routing** so Docker host can reach 192.168.8.0/24
3. **Update `.env`** with appropriate gateway
4. **Configure MikroTik** for monitoring (see [MikroTik Integration](MikroTik-Integration))

---

## Troubleshooting Installation

### Container Won't Start

Check logs for errors:
```bash
docker logs phev2mqtt
```

Common issues:
- Missing or invalid MQTT credentials
- Network connectivity issues
- Permission problems with config directory

### Can't Connect to MQTT Broker

Verify MQTT broker settings:
```bash
# Test MQTT connection
docker exec -it phev2mqtt ping <mqtt_broker_ip>
```

Ensure:
- MQTT broker is running
- Firewall rules allow connection
- Credentials are correct
- MQTT server URL format is correct (e.g., `tcp://192.168.1.2:1883`)

### No Route to PHEV

Check routing:
```bash
docker exec -it phev2mqtt ip route
```

Verify:
- `route_add` is set correctly
- Container has `NET_ADMIN` capability
- Gateway is reachable from Docker host

For more help, see [Troubleshooting](Troubleshooting).

---

## Next Steps

- ✅ [Configure advanced settings](Configuration)
- ✅ [Set up Home Assistant integration](Home-Assistant-Integration)
- ✅ [Configure WiFi management](WiFi-Management)
- ✅ [Review security best practices](Security-Best-Practices)
