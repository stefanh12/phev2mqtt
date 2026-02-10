# Troubleshooting Guide

Common issues and solutions for phev2mqtt.

## Quick Diagnostics

### Check Container Status

```bash
# Docker Compose
docker-compose ps
docker-compose logs -f

# Docker
docker ps | grep phev2mqtt
docker logs phev2mqtt -f

# Unraid
# Navigate to Docker tab, check phev2mqtt status
# Click "Logs" button to view container logs
```

### Check MQTT Connection

```bash
# Test MQTT broker connectivity
mosquitto_pub -h 192.168.1.2 -u user -P pass -t "test" -m "hello"

# Subscribe to phev2mqtt messages
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "phev/#" -v
```

### Enable Debug Logging

```bash
# In .env file
log_level=debug
```

Changes apply within 5 seconds (hot-reload).

---

## Installation Issues

### Container Won't Start

**Symptoms:**
- Container exits immediately after starting
- `docker ps` doesn't show phev2mqtt

**Solutions:**

1. **Check logs for errors:**
   ```bash
   docker logs phev2mqtt
   ```

2. **Common causes:**
   - Missing `.env` file
   - Invalid MQTT credentials
   - Configuration validation errors
   - Port conflicts

3. **Verify .env file exists:**
   ```bash
   ls -la /path/to/phev2mqtt/.env
   # Or for Unraid:
   ls -la /mnt/user/appdata/phev2mqtt/.env
   ```

4. **Test with minimal config:**
   ```bash
   # Minimal .env
   mqtt_server=tcp://192.168.1.2:1883
   mqtt_username=phevmqttuser
   mqtt_password=your_password
   ```

### Permission Denied Errors

**Symptoms:**
- Can't read `.env` file
- Can't add routes
- Network errors

**Solutions:**

1. **Check file permissions:**
   ```bash
   ls -la /path/to/.env
   chmod 644 /path/to/.env
   ```

2. **Verify NET_ADMIN capability:**
   ```yaml
   # In docker-compose.yml
   cap_add:
     - NET_ADMIN
   ```

3. **Check SELinux (Linux):**
   ```bash
   # Temporarily disable for testing
   sudo setenforce 0
   ```

---

## Connection Issues

### Can't Connect to MQTT Broker

**Symptoms:**
- Logs show "Failed to connect to MQTT"
- No messages published to MQTT

**Solutions:**

1. **Verify MQTT broker is reachable:**
   ```bash
   # From Docker host
   ping 192.168.1.2
   
   # From inside container
   docker exec phev2mqtt ping -c 3 192.168.1.2
   ```

2. **Check MQTT server URL format:**
   ```bash
   # Correct formats:
   mqtt_server=tcp://192.168.1.2:1883
   mqtt_server=ssl://mqtt.example.com:8883
   
   # Incorrect (will fail):
   mqtt_server=192.168.1.2:1883  # Missing protocol
   mqtt_server=mqtt://192.168.1.2:1883  # Wrong protocol
   ```

3. **Test MQTT credentials:**
   ```bash
   mosquitto_pub -h 192.168.1.2 -u phevmqttuser -P password \
     -t "test" -m "hello"
   ```

4. **Check firewall rules:**
   ```bash
   # On MQTT broker server
   sudo netstat -tlnp | grep 1883
   sudo iptables -L | grep 1883
   ```

5. **Check MQTT broker logs:**
   ```bash
   # Mosquitto example
   tail -f /var/log/mosquitto/mosquitto.log
   ```

### Can't Connect to PHEV

**Symptoms:**
- "Waiting for PHEV connection" in logs
- No vehicle data received
- Connection timeout errors

**Solutions:**

1. **Verify WiFi connection:**
   ```bash
   # Check if connected to PHEV WiFi
   # On host system (adjust for your OS):
   iwconfig  # Linux
   networksetup -getairportnetwork en0  # macOS
   ```

2. **Check PHEV is reachable:**
   ```bash
   # From Docker host
   ping 192.168.8.46
   
   # From container
   docker exec phev2mqtt ping -c 3 192.168.8.46
   ```

3. **Verify routing:**
   ```bash
   # Check route to PHEV network exists
   docker exec phev2mqtt ip route | grep 192.168.8
   
   # Should see something like:
   # 192.168.8.0/24 via 192.168.1.1 dev eth0
   ```

4. **Check route_add configuration:**
   ```bash
   # In .env file
   route_add=192.168.1.1  # Your gateway IP
   ```

5. **Ensure PHEV WiFi is active:**
   - Vehicle must be within WiFi range
   - PHEV WiFi typically stays active for a while after vehicle is turned off
   - Some models may turn off WiFi after extended periods

6. **Check for registration issues:**
   - First-time setup requires registration (see [Registration](#registration-issues))
   - Verify MAC address hasn't changed (if using NAT/bridge)

### Connection Drops Frequently

**Symptoms:**
- Connection works but disconnects often
- "Connection lost" messages in logs
- Intermittent vehicle data

**Solutions:**

1. **Check WiFi signal strength:**
   - Move WiFi bridge closer to parking spot
   - Check for interference (other 2.4GHz devices)
   - Consider MikroTik bridge for stability

2. **Increase timeout values:**
   ```bash
   # In .env file
   phev_tcp_read_timeout=60s
   phev_tcp_write_timeout=30s
   connection_retry_interval=30s
   ```

3. **Enable remote WiFi restart:**
   ```bash
   # In .env file
   remote_wifi_restart_enabled=true
   remote_wifi_restart_topic=mikrotik/phev/restart
   remote_wifi_restart_message=restart
   ```

4. **Check for encoding errors:**
   ```bash
   # Look for XOR/encoding errors in logs
   docker logs phev2mqtt | grep -i "xor\|encoding\|bad"
   ```

---

## Registration Issues

### First-Time Registration Fails

**Symptoms:**
- Can't complete initial vehicle registration
- Registration timeout
- "Registration failed" errors

**Solutions:**

1. **Verify registration mode:**
   ```bash
   # In .env file
   phev_register=true
   ```

2. **Put vehicle in registration mode:**
   - Consult your vehicle manual for exact procedure
   - Usually involves pressing buttons on remote app display
   - Vehicle must be in range and WiFi active

3. **Check registration timeout:**
   ```bash
   # In .env file (increase if needed)
   phev_register_timeout=20s
   ```

4. **Monitor logs during registration:**
   ```bash
   docker logs phev2mqtt -f
   # Look for registration prompts and errors
   ```

5. **Ensure clean state:**
   - Unregister any existing clients if MAX reached
   - Restart container after setting phev_register=true

6. **After successful registration:**
   ```bash
   # In .env file
   phev_register=false
   # Then restart container
   ```

### Already Registered But Not Connecting

**Symptoms:**
- Previously worked, now won't connect
- "Not registered" or "Authentication failed" errors

**Solutions:**

1. **Check if MAC address changed:**
   - Vehicle registers MAC address
   - If using NAT/bridge, MAC must remain consistent
   - Re-register if MAC changed

2. **Verify registration count:**
   - Check register 0x15 for number of registered clients
   - Vehicle may have MAX registrations (typically 4)
   - May need to unregister old clients

3. **Force unregister and re-register:**
   ```bash
   # Use phev2mqtt unregister command
   ./phev2mqtt client register --unregister
   
   # Then register again
   phev_register=true
   docker-compose restart
   ```

---

## Home Assistant Issues

### Entities Not Appearing

**Symptoms:**
- Home Assistant doesn't show PHEV device
- No auto-discovery

**Solutions:**

1. **Check MQTT integration:**
   - Settings → Devices & Services → MQTT
   - Verify "Connected" status
   - Check broker configuration

2. **Check discovery messages:**
   ```bash
   # Subscribe to discovery topic
   mosquitto_sub -h 192.168.1.2 -u user -P pass \
     -t "homeassistant/#" -v
   ```

3. **Check phev2mqtt logs:**
   ```bash
   docker logs phev2mqtt | grep -i "discovery"
   ```

4. **Force discovery republish:**
   ```bash
   # Restart phev2mqtt
   docker-compose restart
   ```

5. **Set VIN for immediate discovery:**
   ```bash
   # In .env file
   vehicle_vin=JA4J24A58KZ123456
   ```

### Entities Show "Unavailable"

**Symptoms:**
- Entities exist but show "Unavailable"
- No data updates

**Solutions:**

1. **Check phev2mqtt is connected to PHEV:**
   ```bash
   docker logs phev2mqtt | grep -i "connected\|phev"
   ```

2. **Check availability topic:**
   ```bash
   mosquitto_sub -h 192.168.1.2 -u user -P pass \
     -t "phev/availability" -v
   # Should show "online"
   ```

3. **Verify phev2mqtt container is running:**
   ```bash
   docker ps | grep phev2mqtt
   ```

4. **Check MQTT connection:**
   - Broker may be down or restarting
   - Network issues between phev2mqtt and broker

### State Updates Not Working

**Symptoms:**
- Can see entities but values don't update
- Old/stale data

**Solutions:**

1. **Check update interval:**
   ```bash
   # In .env file
   update_interval=5m  # Reduce if needed
   ```

2. **Verify PHEV connection:**
   - Vehicle must be in range
   - WiFi must be active
   - Check connection in logs

3. **Check for MQTT retained messages:**
   ```bash
   # Clear all phev topics
   mosquitto_pub -h 192.168.1.2 -u user -P pass \
     -t "phev/#" -r -n
   ```

4. **Restart phev2mqtt:**
   ```bash
   docker-compose restart
   ```

---

## MikroTik Integration Issues

### MikroTik Won't Connect to PHEV WiFi

**Symptoms:**
- MikroTik WiFi interface shows disconnected
- Can't get DHCP lease from PHEV

**Solutions:**

1. **Check WiFi configuration:**
   ```routeros
   /interface wireless print
   /interface wireless security-profiles print
   ```

2. **Verify SSID and password:**
   - Must match PHEV WiFi exactly
   - Check for typos in security profile

3. **Check frequency:**
   ```routeros
   /interface wireless set wlan1 frequency=2422
   ```

4. **Enable interface:**
   ```routeros
   /interface wireless enable wlan1
   ```

5. **Check logs:**
   ```routeros
   /log print where topics~"wireless"
   ```

### MQTT Not Publishing from MikroTik

**Symptoms:**
- No MQTT messages from MikroTik
- Connection status not updating

**Solutions:**

1. **Check MQTT broker configuration:**
   ```routeros
   /iot mqtt brokers print
   # Should show "Connected"
   ```

2. **Verify broker credentials:**
   ```routeros
   /iot mqtt brokers set homeassistantmqtt \
     username=user password=pass
   ```

3. **Test manual publish:**
   ```routeros
   /iot mqtt publish broker=homeassistantmqtt \
     topic=test message="hello"
   ```

4. **Check routing to broker:**
   ```routeros
   /ping 192.168.1.197 count=3
   /ip route print
   ```

5. **Check MQTT logs:**
   ```routeros
   /log print where message~"mqtt"
   ```

### WiFi Control Not Working

**Symptoms:**
- Can't Turn WiFi on/off via MQTT
- Power save mode not working

**Solutions:**

1. **Check MQTT subscription:**
   ```routeros
   /iot mqtt subscriptions print
   ```

2. **Test subscription manually:**
   ```bash
   # Publish test message
   mosquitto_pub -h 192.168.1.197 -u user -P pass \
     -t "homeassistant/sensor/mikrotik/wifi" \
     -m '{"wifi": "disable"}'
   ```

3. **Check subscription script:**
   - Verify JSON format matches subscription on-message script
   - Check for escaping issues in RouterOS script

4. **Monitor RouterOS logs:**
   ```routeros
   /log print follow where message~"Got data"
   ```

5. **Verify phev2mqtt configuration:**
   ```bash
   # In .env file
   remote_wifi_control_topic=homeassistant/sensor/mikrotik/wifi
   remote_wifi_enable_message={"wifi": "enable"}
   remote_wifi_disable_message={"wifi": "disable"}
   ```

---

## Configuration Issues

### Hot-Reload Not Working

**Symptoms:**
- Changes to `.env` file not applying
- Must restart container for changes

**Solutions:**

1. **Check which settings changed:**
   - Some settings require restart (MQTT credentials, phev_register)
   - Most timeouts and intervals are hot-reloadable

2. **Verify file save:**
   ```bash
   # Check file timestamp
   ls -la /path/to/.env
   ```

3. **Check logs for reload messages:**
   ```bash
   docker logs phev2mqtt | grep -i "reload"
   # Should see: "Configuration file changed, reloading..."
   ```

4. **Check file permissions:**
   ```bash
   chmod 644 /path/to/.env
   ```

5. **Settings requiring restart:**
   - mqtt_server
   - mqtt_username
   - mqtt_password
   - mqtt_topic_prefix
   - phev_register
   - route_add

### Configuration Validation Errors

**Symptoms:**
- Container won't start
- "Invalid configuration" errors in logs

**Solutions:**

1. **Check logs for specific error:**
   ```bash
   docker logs phev2mqtt 2>&1 | grep -i "error\|invalid"
   ```

2. **Common validation errors:**
   - MQTT password too short (min 8 chars, recommended 12+)
   - Invalid MQTT server URL format
   - Invalid topic names (containing wildcards +, #)
   - Invalid duration format (use s, m, h units)

3. **Verify MQTT server format:**
   ```bash
   # Correct:
   mqtt_server=tcp://192.168.1.2:1883
   mqtt_server=ssl://mqtt.example.com:8883
   
   # Wrong:
   mqtt_server=192.168.1.2  # Missing protocol and port
   ```

4. **Check command injection in wifi_restart_command:**
   - No semicolons (;)
   - No pipes (|)
   - No command substitution ($(), backticks)
   - && allowed but warned

---

## Performance Issues

### High CPU Usage

**Symptoms:**
- Container using excessive CPU
- System slow when phev2mqtt running

**Solutions:**

1. **Check for connection loops:**
   ```bash
   docker logs phev2mqtt | grep -i "retry\|reconnect"
   ```

2. **Increase retry intervals:**
   ```bash
   # In .env file
   retry_interval=120s
   connection_retry_interval=120s
   ```

3. **Reduce logging:**
   ```bash
   # In .env file
   log_level=info  # Or warning/error
   ```

4. **Check for encoding errors:**
   ```bash
   docker logs phev2mqtt | grep -i "xor\|error"
   ```

### High Memory Usage

**Symptoms:**
- Container using lots of RAM
- Memory growing over time

**Solutions:**

1. **Restart container periodically:**
   ```bash
   docker-compose restart
   ```

2. **Check for memory leaks:**
   ```bash
   docker stats phev2mqtt
   ```

3. **Report issue:**
   - If memory grows continuously, report on GitHub with logs

---

## Debugging Tools

### Enable Verbose Logging

```bash
# In .env file
log_level=debug
```

### Decode Hex Packets

```bash
./phev2mqtt decode <hex_packet>
```

### Analyze Packet Capture

```bash
./phev2mqtt pcap <capture_file.pcap>
```

### Monitor MQTT Topics

```bash
# All phev2mqtt messages
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "phev/#" -v

# MikroTik messages
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "mikrotik/#" -v

# Home Assistant discovery
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "homeassistant/#" -v
```

### Container Shell Access

```bash
# Access container shell
docker exec -it phev2mqtt /bin/sh

# Check routes
ip route

# Check network
ping 192.168.8.46

# Check DNS
nslookup mqtt.example.com
```

---

## Getting Help

If you're still experiencing issues:

1. **Gather information:**
   - Full container logs: `docker logs phev2mqtt > logs.txt`
   - Configuration (remove passwords): `cat .env`
   - System info: Docker version, host OS, network setup

2. **Search existing issues:**
   - [GitHub Issues](https://github.com/stefanh12/phev2mqtt/issues)
   - Check closed issues too

3. **Create new issue:**
   - Provide logs and configuration
   - Describe expected vs actual behavior
   - List troubleshooting steps already tried

4. **Community help:**
   - [GitHub Discussions](https://github.com/stefanh12/phev2mqtt/discussions)

---

## Next Steps

- [Configuration](Configuration) - Review configuration options
- [Security Best Practices](Security-Best-Practices) - Secure your installation
- [Development](Development) - Contribute fixes or enhancements
