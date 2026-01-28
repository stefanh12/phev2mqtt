# phev2mqtt Unraid Setup

This folder contains the necessary files to run phev2mqtt in a Docker container on Unraid, as an alternative to running it on Home Assistant.

## 📦 Contents

- `XMLFile1.xml` - Community Applications (CA) template for Unraid
- `Dockerfile` - Docker image build configuration
- `entrypoint.sh` - Container startup script

## 🚀 Installation

### Quick Install via Community Applications

1. Copy `XMLFile1.xml` to `/boot/config/plugins/dockerMan/templates-user` on your Unraid USB drive
2. In Unraid's Docker tab, click **Add Container** → **Select Template**
3. Find and select **phev2mqtt** from the templates list
4. Configure the required parameters (see Configuration section below)
5. Click **Apply** to create and start the container

### Manual Docker Hub Installation

The pre-built Docker image is available at: https://hub.docker.com/r/hstefan/phev/

## ⚙️ Configuration

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `mqtt_server` | MQTT broker address and port | `192.168.1.2:1883` |
| `mqtt_user` | MQTT username | `phevmqttuser` |
| `mqtt_password` | MQTT password | `mqttuserpassword` |

### Optional Environment Variables

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `debug` | Enable debug mode (sleeps container indefinitely for troubleshooting) | `false` | `true`/`false` |
| `phev_register` | Register client with PHEV (must be set to `true` for initial setup) | `false` | `true`/`false` |
| `route_add` | Gateway IP that routes traffic to your PHEV | - | IP address (e.g., `192.168.1.1`) |

### Network Configuration

- **Network Mode**: Custom bridge (e.g., `br0.308`)
- **IP Address**: Static IP recommended (e.g., `192.168.7.10`)
- **Extra Parameters**: `--cap-add=NET_ADMIN` (required for network routing)

### Routing Setup

If your PHEV is on a different network segment (192.168.8.0/24), you need to set the `route_add` variable to your gateway IP:

```bash
route_add=192.168.1.1
```

This adds the route: `192.168.8.0/24 via <gateway> dev eth0`

## 🔧 Initial Setup - Client Registration

Before first use, you must register the client with your PHEV:

1. Set `phev_register=true` in the container configuration
2. Start the container - it will run the registration process
3. Follow any on-screen prompts in the container logs
4. After successful registration, set `phev_register=false` and restart

## 🐛 Debugging

To troubleshoot connection or configuration issues:

1. Set `debug=true` in the container configuration
2. The container will sleep indefinitely
3. Access the container console via Unraid
4. Manually run registration or test commands:
   ```bash
   /usr/src/app/phev2mqtt/phev2mqtt client register
   ```

## ✅ Tested Configuration

This setup has been tested with:
- **Home Assistant**: Running on Raspberry Pi 4
- **Unraid**: Server for running phev2mqtt container
- **WiFi Bridge**: Mikrotik SXTsq Lite2

### Mikrotik SXTsq Configuration

The SXTsq requires the **MQTT package add-on** to be installed. This enables Home Assistant to monitor the WiFi connection status and automate actions.

**Common Issue**: If the PHEV loses WiFi connection and the link goes down, both the SXTsq and phev2mqtt may need to be restarted (not just the SXTsq).

## 📚 Additional Resources

- [phev2mqtt GitHub Repository](https://github.com/stefanh12/phev2mqtt)
- [Docker Hub Image](https://hub.docker.com/r/hstefan/phev/)
- Main project README for more details about phev2mqtt functionality

## 🤝 Support

For issues or questions, please open an issue on the [GitHub repository](https://github.com/stefanh12/phev2mqtt/issues).


