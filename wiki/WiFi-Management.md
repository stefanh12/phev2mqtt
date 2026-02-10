# WiFi Management Guide

phev2mqtt includes advanced WiFi management features to optimize connection reliability and reduce power consumption.

## Overview

WiFi management helps with:
- **Automatic recovery** from connection issues
- **Power saving** by turning off WiFi when not needed
- **Remote control** of WiFi devices via MQTT
- **Stability improvements** with external WiFi bridges

## WiFi Management Features

### 1. Local WiFi Restart

Automatically restart the local WiFi interface when connection is lost.

**Use case:** phev2mqtt running on hardware with its own WiFi interface.

**Configuration:**
```bash
# In .env file
local_wifi_restart_enabled=true
wifi_restart_time=10m
wifi_restart_command=nmcli device disconnect wlan0 && nmcli device connect wlan0
```

**How it works:**
- Monitors connection to PHEV
- If connection lost for `wifi_restart_time`, restarts WiFi
- Executes custom command if provided

**Note:** Only works on some hardware configurations. Test before deploying.

### 2. Remote WiFi Restart (MikroTik)

Send MQTT commands to remotely restart WiFi on external devices.

**Use case:** Using a dedicated WiFi bridge (e.g., MikroTik RBSXTsq2nD).

**Configuration:**
```bash
# In .env file
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart
remote_wifi_restart_min_interval=2m
```

**How it works:**
- Monitors connection to PHEV
- When connection lost, publishes MQTT message to restart topic
- External device (MikroTik) receives message and restarts WiFi
- Minimum interval prevents excessive restarts

**Setup required:**
- Configure MikroTik MQTT subscription (see [MikroTik Integration](MikroTik-Integration))
- MQTT broker must be reachable by both phev2mqtt and MikroTik

### 3. Remote WiFi Power Save Mode

Automatically turn off WiFi between update intervals to save power.

**Use case:** Reduce power consumption on WiFi bridge when not communicating with PHEV.

**Configuration:**
```bash
# In .env file
remote_wifi_power_save_enabled=true
remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi
remote_wifi_enable_message={"wifi": "enable"}
remote_wifi_disable_message={"wifi": "disable"}
remote_wifi_power_save_wait=5s
remote_wifi_command_wait=10s

# Required: Update interval must be > 1 minute
update_interval=10m
```

**How it works:**
- Before each update, publishes "enable" message
- Waits `remote_wifi_power_save_wait` for WiFi to establish link
- Performs update
- After update, publishes "disable" message
- Keeps WiFi on longer (`remote_wifi_command_wait`) after manual commands

**Requirements:**
- `update_interval` must be greater than 1 minute
- `remote_wifi_control_topic` must be configured
- External device must support WiFi on/off via MQTT

**Benefits:**
- Reduces power consumption on WiFi bridge
- Extends battery life for battery-powered bridges
- Reduces RF exposure

---

## Choosing the Right Approach

### Local WiFi Restart

**✅ Good for:**
- Single-board computers (Raspberry Pi, etc.)
- Systems with built-in WiFi
- Simple setups

**❌ Not ideal for:**
- Containers without host network access
- Systems without WiFi management tools
- Production environments (unreliable)

### Remote WiFi Restart

**✅ Good for:**
- Dedicated WiFi bridges
- MikroTik devices
- Stable, monitored deployments
- Separate VLANs

**Recommended hardware:**
- MikroTik RBSXTsq2nD (outdoor, 2.4GHz)
- Any MikroTik with 2.4GHz WiFi and MQTT support

### Remote Power Save Mode

**✅ Good for:**
- Battery-powered bridges
- Energy-conscious deployments
- Infrequent updates (10+ minutes)

**❌ Not ideal for:**
- Frequent updates (< 5 minutes)
- Real-time monitoring requirements
- Immediate command response needed

---

## Configuration Details

### Local WiFi Restart Settings

**local_wifi_restart_enabled**
- Enable local WiFi restart feature
- Default: `false`

**wifi_restart_time**
- Duration without connection before restarting
- Default: `10m`
- Format: Duration with units (s, m, h)

**wifi_restart_command**
- Custom command to restart WiFi
- Default: Empty (uses system default)
- Security: Validated to prevent command injection
- Example: `nmcli device disconnect wlan0 && nmcli device connect wlan0`

### Remote WiFi Restart Settings

**remote_wifi_restart_enabled**
- Enable remote restart via MQTT
- Default: `false`

**remote_wifi_restart_topic**
- MQTT topic to publish restart command
- Default: `mikrotik/phev/restart`
- Validated for security

**remote_wifi_restart_message**
- Restart command message payload
- Default: `restart`

**remote_wifi_restart_min_interval**
- Minimum time between restart attempts
- Default: `2m`
- Prevents excessive restarts

### Remote Power Save Settings

**remote_wifi_power_save_enabled**
- Enable power save mode
- Default: `false`

**remote_wifi_control_topic**
- MQTT topic for on/off commands
- Required when power save enabled
- Example: `homeassistant/sensor/mikrotik/wifi`

**remote_wifi_enable_message**
- Message to turn WiFi on
- Default: `{"wifi": "enable"}`

**remote_wifi_disable_message**
- Message to turn WiFi off
- Default: `{"wifi": "disable"}`

**remote_wifi_power_save_wait**
- Wait time after enabling WiFi
- Default: `5s`
- Allows link establishment

**remote_wifi_command_wait**
- Keep WiFi on after commands
- Default: `10s`
- Ensures status updates received

---

## Example Configurations

### Simple Remote Restart (MikroTik)

```bash
# MQTT Settings
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!

# Basic settings
log_level=info
update_interval=5m

# Remote WiFi restart only
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart
```

### Full Power Save Setup

```bash
# MQTT Settings
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!

# Basic settings
log_level=info
update_interval=10m  # Must be > 1 minute for power save

# Remote restart
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart
remote_wifi_restart_min_interval=2m

# Power save mode
remote_wifi_power_save_enabled=true
remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi
remote_wifi_enable_message={"wifi": "enable"}
remote_wifi_disable_message={"wifi": "disable"}
remote_wifi_power_save_wait=5s
remote_wifi_command_wait=10s
```

### Local WiFi Restart (Raspberry Pi)

```bash
# MQTT Settings
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=SecurePassword123!

# Basic settings
log_level=info
update_interval=5m

# Local WiFi restart
local_wifi_restart_enabled=true
wifi_restart_time=10m
wifi_restart_command=sudo systemctl restart wpa_supplicant
```

---

## Monitoring and Debugging

### Check WiFi Management in Logs

```bash
docker logs phev2mqtt | grep -i "wifi"
```

Look for:
- "Enabling remote WiFi" / "Disabling remote WiFi"
- "Publishing WiFi restart command"
- "Restarting local WiFi interface"

### Monitor MQTT Messages

```bash
# Subscribe to control topic
mosquitto_sub -h 192.168.1.2 -u user -P pass \
  -t "homeassistant/sensor/mikrotik/wifi" -v

# Subscribe to restart topic
mosquitto_sub -h 192.168.1.2 -u user -P pass \
  -t "mikrotik/phev/restart" -v
```

### Test Power Save Manually

```bash
# Enable WiFi
mosquitto_pub -h 192.168.1.2 -u user -P pass \
  -t "homeassistant/sensor/mikrotik/wifi" \
  -m '{"wifi": "enable"}'

# Wait a moment, then disable
mosquitto_pub -h 192.168.1.2 -u user -P pass \
  -t "homeassistant/sensor/mikrotik/wifi" \
  -m '{"wifi": "disable"}'
```

---

## Best Practices

### For Stability

✅ Use dedicated WiFi bridge (MikroTik)  
✅ Position bridge close to parking spot  
✅ Use 2.4GHz (better range than 5GHz)  
✅ Monitor connection status via Home Assistant  
✅ Set appropriate restart intervals  

### For Power Saving

✅ Use power save mode with longer update intervals (10-15 minutes)  
✅ Balance update frequency with battery life  
✅ Monitor WiFi on/off cycles  
✅ Adjust wait times based on your WiFi device  

### For Reliability

✅ Test all WiFi management features before production use  
✅ Monitor logs during initial setup  
✅ Set reasonable timeout values  
✅ Use MikroTik for best stability (MY20+ PHEV)  
✅ Keep firmware updated  

---

## Troubleshooting

### Power Save Not Working

**Check:**
- `update_interval` is > 1 minute
- `remote_wifi_control_topic` is configured
- External device responds to on/off commands
- MQTT messages are being published

**Test manually:**
```bash
mosquitto_sub -h 192.168.1.2 -u user -P pass \
  -t "homeassistant/sensor/mikrotik/wifi" -v
```

### Remote Restart Not Triggering

**Check:**
- `remote_wifi_restart_enabled=true`
- Connection actually lost (check logs)
- MQTT broker reachable
- External device subscribed to topic
- Minimum interval not blocking restart

**Test manually:**
```bash
mosquitto_pub -h 192.168.1.2 -u user -P pass \
  -t "mikrotik/phev/restart" -m "restart"
```

### Local Restart Fails

**Check:**
- Command has proper permissions
- Command syntax is correct
- Network manager tool installed (`nmcli`, etc.)
- Container has necessary capabilities

**Test command manually:**
```bash
docker exec phev2mqtt /bin/sh -c "your_wifi_command_here"
```

See [Troubleshooting](Troubleshooting) for more solutions.

---

## Next Steps

- [MikroTik Integration](MikroTik-Integration) - Set up dedicated WiFi bridge
- [Configuration](Configuration) - Detailed configuration options
- [Home Assistant Integration](Home-Assistant-Integration) - Monitor WiFi status
- [Troubleshooting](Troubleshooting) - Common WiFi issues
