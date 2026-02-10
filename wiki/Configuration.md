# Configuration Guide

Complete reference for configuring phev2mqtt through the `.env` file or environment variables.

## Configuration File Location

phev2mqtt looks for configuration in the following order:

1. `/config/.env` (Docker mount point)
2. `.env` in the current directory
3. Environment variables

## Hot Reload Support

The application automatically watches the `.env` file for changes and reloads configuration **without requiring a container restart**. Changes are detected every 5 seconds.

### Hot-Reloadable Settings

These settings can be changed while the container is running:
- `update_interval` - Update request frequency
- `retry_interval` - Connection retry frequency
- `log_level` - Logging verbosity
- `remote_wifi_restart_topic` - Remote restart MQTT topic
- `remote_wifi_restart_message` - Remote restart message
- `remote_wifi_control_topic` - WiFi power control topic
- `remote_wifi_enable_message` - WiFi enable message
- `remote_wifi_disable_message` - WiFi disable message
- All timeout settings
- Remote WiFi control and power save settings

### Settings Requiring Restart

These settings require a container restart to take effect:
- `mqtt_server` - MQTT broker address
- `mqtt_username` - MQTT username
- `mqtt_password` - MQTT password
- `mqtt_topic_prefix` - MQTT topic prefix
- `phev_register` - Registration mode
- `route_add` - Network routing configuration

### Applying Changes

Simply edit your `.env` file and save. Changes will be applied within 5 seconds. Check the logs to confirm:
```
Configuration file changed, reloading...
Configuration reloaded successfully
```

---

## Required Settings

### MQTT Configuration

**mqtt_server** (Required)  
MQTT broker address with protocol and port.

- **Format**: `protocol://host:port`
- **Protocols**: `tcp://`, `ssl://`, `tls://`, `ws://`, `wss://`
- **Examples**:
  - `tcp://192.168.1.2:1883` (unencrypted)
  - `ssl://mqtt.example.com:8883` (TLS encrypted)
  - `wss://mqtt.example.com:8084` (WebSocket TLS)
- **Security Warning**: Unencrypted connections (`tcp://`, `ws://`) will show a warning. Use TLS in production.

**mqtt_username** (Required)  
MQTT broker username for authentication.

**mqtt_password** (Required)  
MQTT broker password.

- **Minimum**: 8 characters (hard requirement)
- **Recommended**: 12+ characters
- **Validation**: Blocks common weak passwords
- **Warning**: Password is stored in plaintext in `.env` file. Protect this file!

---

## Basic Settings

### Logging

**log_level**  
Control logging verbosity.

- **Default**: `info`
- **Options**:
  - `none` - Only fatal errors
  - `error` - Error messages only
  - `warning` - Warnings and errors
  - `info` - Normal operation info (recommended)
  - `debug` - Detailed debug information
- **Example**: `log_level=info`

### Update Intervals

**update_interval**  
How often to request updates from the PHEV.

- **Default**: `5m`
- **Format**: Duration with units (`s`, `m`, `h`)
- **Examples**: `5m`, `10m`, `15m`
- **Note**: Shorter intervals consume more vehicle battery
- **Hot-reloadable**: Yes

**retry_interval**  
How often to retry when connection to PHEV is lost.

- **Default**: `60s`
- **Format**: Duration with units (`s`, `m`, `h`)
- **Examples**: `30s`, `1m`, `2m`
- **Hot-reloadable**: Yes

### MQTT Topic Prefix

**mqtt_topic_prefix**  
Base topic prefix for all MQTT messages.

- **Default**: `phev`
- **Example**: `phev` results in topics like `phev/battery`, `phev/temperature`
- **Validation**: Blocks MQTT wildcards (+, #), control characters
- **Requires restart**: Yes

### Vehicle Identification

**vehicle_vin**  
Your vehicle's VIN (Vehicle Identification Number).

- **Default**: Empty (VIN read from vehicle)
- **Optional**: If set, enables immediate Home Assistant discovery on startup
- **Behavior**: 
  - **Empty**: Discovery waits for first PHEV connection
  - **Set**: Discovery happens immediately on startup
- **Example**: `vehicle_vin=JA4J24A58KZ123456`

---

## Network Configuration

### Routing

**route_add**  
Gateway IP address for routing to the PHEV network (192.168.8.0/24).

- **Default**: Empty (no route added)
- **Example**: `route_add=192.168.1.1`
- **Usage**: Required when Docker host needs routing to reach PHEV
- **Note**: Container requires `NET_ADMIN` capability
- **Requires restart**: Yes

---

## PHEV Settings

### Registration Mode

**phev_register**  
Enable first-time vehicle registration mode.

- **Default**: `false`
- **Values**: `true` or `false`
- **Usage**:
  1. Set to `true` for first-time setup
  2. Follow prompts in logs to enter security code from vehicle dashboard
  3. After successful registration, set to `false`
  4. Restart container
- **Warning**: Only use during initial registration. Keep `false` for normal operation.
- **Requires restart**: Yes

---

## WiFi Management

### Local WiFi Restart

Automatically restart the local WiFi interface when connection to PHEV is lost.

**local_wifi_restart_enabled**  
Enable local WiFi restart feature.

- **Default**: `false`
- **Values**: `true` or `false`
- **Note**: Only works on some hardware configurations

**wifi_restart_time**  
Duration without connection before restarting WiFi.

- **Default**: `10m`
- **Format**: Duration with units
- **Example**: `wifi_restart_time=10m`

**wifi_restart_command**  
Custom command to restart WiFi interface.

- **Default**: Empty (uses default system command)
- **Example**: `wifi_restart_command=nmcli device disconnect wlan0 && nmcli device connect wlan0`
- **Security**: Validated to prevent command injection (blocks `;`, `|`, `$()`, backticks)
- **Allowed**: `&&` for command chaining (with warning)

### Remote WiFi Restart (MikroTik)

Send MQTT commands to remotely restart WiFi on external devices when connection is lost.

**remote_wifi_restart_enabled**  
Enable remote WiFi restart via MQTT.

- **Default**: `false`
- **Values**: `true` or `false`
- **Use Case**: MikroTik WiFi bridge or access point

**remote_wifi_restart_topic**  
MQTT topic to publish restart command to.

- **Default**: `mikrotik/phev/restart`
- **Example**: `remote_wifi_restart_topic=mikrotik/phev/restart`
- **Validation**: Topic format validated for security
- **Hot-reloadable**: Yes

**remote_wifi_restart_message**  
Message payload to send for restart.

- **Default**: `restart`
- **Example**: `remote_wifi_restart_message=restart`
- **Hot-reloadable**: Yes

**remote_wifi_restart_min_interval**  
Minimum time between restart attempts.

- **Default**: `2m`
- **Format**: Duration with units
- **Purpose**: Prevents excessive restart commands

See [MikroTik Integration](MikroTik-Integration) for complete setup instructions.

### Remote WiFi Power Save Mode

Automatically turn off WiFi between update intervals to save power.

**remote_wifi_power_save_enabled**  
Enable power save mode.

- **Default**: `false`
- **Values**: `true` or `false`
- **Requirements**:
  - `update_interval` must be > 1 minute
  - `remote_wifi_control_topic` must be configured
- **Behavior**:
  - WiFi turned on before each update
  - WiFi turned off after update completion
  - WiFi stays on longer after manual commands (to receive status)

**remote_wifi_control_topic**  
MQTT topic for WiFi control commands.

- **Required**: When power save is enabled
- **Example**: `remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi`
- **Validation**: Topic format validated
- **Hot-reloadable**: Yes

**remote_wifi_enable_message**  
Message to enable WiFi.

- **Default**: `{"wifi": "enable"}`
- **Format**: JSON or plain text
- **Example**: `remote_wifi_enable_message={"wifi": "enable"}`
- **Hot-reloadable**: Yes

**remote_wifi_disable_message**  
Message to disable WiFi.

- **Default**: `{"wifi": "disable"}`
- **Format**: JSON or plain text
- **Example**: `remote_wifi_disable_message={"wifi": "disable"}`
- **Hot-reloadable**: Yes

**remote_wifi_power_save_wait**  
Time to wait after enabling WiFi for link establishment.

- **Default**: `5s`
- **Format**: Duration with units
- **Purpose**: Allow time for WiFi connection to stabilize
- **Hot-reloadable**: Yes

**remote_wifi_command_wait**  
Time to keep WiFi on after commands to receive status updates.

- **Default**: `10s`
- **Format**: Duration with units
- **Purpose**: Ensure status updates are received after commands
- **Hot-reloadable**: Yes

---

## Advanced Timeout Configuration

⚠️ **Warning**: Only modify these if you understand their impact. Incorrect values may cause connection issues, increased battery drain, or unexpected behavior.

### Connection Timeouts

**connection_retry_interval**  
Time between connection retry attempts.

- **Default**: `60s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

**availability_offline_timeout**  
Time before publishing MQTT offline status.

- **Default**: `30s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

### PHEV Communication Timeouts

**phev_start_timeout**  
Timeout for PHEV start command response.

- **Default**: `20s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

**phev_register_timeout**  
Timeout for register set acknowledgment.

- **Default**: `10s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

**phev_tcp_read_timeout**  
TCP read deadline for PHEV connection.

- **Default**: `30s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

**phev_tcp_write_timeout**  
TCP write deadline for PHEV connection.

- **Default**: `15s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

### Other Timeouts

**encoding_error_reset_interval**  
Time after which encoding error count resets.

- **Default**: `15s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

**config_reload_interval**  
How often to check for configuration file changes.

- **Default**: `5s`
- **Format**: Duration with units
- **Hot-reloadable**: Yes

---

## Complete Configuration Examples

### Minimal Configuration

```bash
# Required only
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_secure_password
```

### Basic Production Configuration

```bash
# MQTT Broker
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_secure_password

# Network
route_add=192.168.1.1

# Logging
log_level=info

# Update Intervals
update_interval=5m
retry_interval=60s

# Vehicle
vehicle_vin=JA4J24A58KZ123456
phev_register=false
```

### Advanced Configuration with WiFi Management

```bash
# MQTT Broker (TLS encrypted)
mqtt_server=ssl://mqtt.example.com:8883
mqtt_username=phevmqttuser
mqtt_password=your_very_secure_password_here

# Network
route_add=192.168.1.1

# Logging
log_level=info

# Update Intervals
update_interval=10m
retry_interval=60s

# Vehicle
vehicle_vin=JA4J24A58KZ123456
phev_register=false

# Remote WiFi Restart (MikroTik)
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart
remote_wifi_restart_min_interval=2m

# Remote WiFi Power Save
remote_wifi_power_save_enabled=true
remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi
remote_wifi_enable_message={"wifi": "enable"}
remote_wifi_disable_message={"wifi": "disable"}
remote_wifi_power_save_wait=5s
remote_wifi_command_wait=10s
```

### Debug Configuration

```bash
# MQTT Broker
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_secure_password

# Network
route_add=192.168.1.1

# Enable debug logging
log_level=debug

# Longer update interval during debugging
update_interval=15m
retry_interval=30s

# Vehicle
phev_register=false

# Adjust timeouts for debugging
phev_tcp_read_timeout=60s
phev_tcp_write_timeout=30s
```

---

## Environment Variables

All `.env` settings can also be set as environment variables using UPPERCASE names:

```bash
MQTT_SERVER=tcp://192.168.1.2:1883
MQTT_USERNAME=phevmqttuser
MQTT_PASSWORD=your_secure_password
LOG_LEVEL=info
UPDATE_INTERVAL=5m
```

**Precedence**: `.env` file settings take precedence over environment variables.

---

## Security Considerations

See [Security Best Practices](Security-Best-Practices) for comprehensive security guidance.

**Quick tips:**
- ✅ Use strong MQTT passwords (12+ characters)
- ✅ Use TLS for MQTT connections (`ssl://` or `tls://`)
- ✅ Protect your `.env` file (contains credentials)
- ✅ Regularly review logs for suspicious activity
- ✅ Keep container images updated
- ✅ Use separate VLANs for PHEV network

---

## Next Steps

- [Home Assistant Integration](Home-Assistant-Integration) - Set up auto-discovery
- [WiFi Management](WiFi-Management) - Optimize WiFi connection
- [MikroTik Integration](MikroTik-Integration) - Configure RouterOS
- [Troubleshooting](Troubleshooting) - Common issues and solutions
