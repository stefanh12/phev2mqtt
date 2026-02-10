# MikroTik Integration

This guide explains how to set up a MikroTik WiFi bridge (specifically the RBSXTsq2nD) to provide a stable, monitored connection to your Mitsubishi Outlander PHEV.

## Why Use a MikroTik Bridge?

✅ **Highly Stable** - Tested with MY20+ PHEV models  
✅ **Automatic Monitoring** - Ping-based connection checks  
✅ **MQTT Integration** - Publishes status updates to Home Assistant  
✅ **Remote Management** - Control WiFi via MQTT commands  
✅ **Power Saving** - Turn WiFi on/off automatically  
✅ **Separate VLAN** - Isolate PHEV network from main network  

## Recommended Hardware

**MikroTik RBSXTsq2nD**
- 2.4GHz WiFi client
- Outdoor rated
- Can be mounted in garage near parking spot
- Highly stable connection to PHEV
- Supports RouterOS 7.x with MQTT

## Architecture Overview

```
PHEV WiFi (192.168.8.46)
    ↓
MikroTik RBSXTsq2nD (WiFi Client)
    ↓
Home Network (VLAN 308)
    ↓
Unraid/Docker Host
    ↓
phev2mqtt Container
    ↓
MQTT Broker
    ↓
Home Assistant
```

---

## Complete RouterOS Configuration

Below is a complete working configuration for the RBSXTsq2nD. This configuration includes:

- WiFi client connecting to PHEV
- VLAN interface for home network
- MQTT broker connection
- DHCP script to publish WiFi connection status
- Ping monitoring script (optional, can be disabled)
- Status publishing to Home Assistant

### Prerequisites

1. **MikroTik RBSXTsq2nD** with RouterOS 7.21.2 or newer
2. **MQTT Broker** (e.g., Mosquitto) on your network
3. **Home Assistant** with MQTT integration
4. **Network Access** - MikroTik must reach MQTT broker
5. **PHEV WiFi Credentials** - SSID and password from your vehicle

### Full Configuration Example

See the complete configuration in [routeros.md](../routeros.md) in the repository.

**Key Configuration Sections:**

1. **WiFi Client Interface**
2. **VLAN Interface for Home Network**
3. **MQTT Broker Configuration**
4. **MQTT Subscriptions for WiFi Control**
5. **DHCP Script for Connection Status**
6. **Scheduler and Scripts**

---

## Step-by-Step Setup

### Step 1: Basic RouterOS Setup

1. **Connect to MikroTik** via WinBox or web interface
2. **Set identity:**
   ```routeros
   /system identity set name=SZTsqlite2garage
   ```
3. **Set timezone:**
   ```routeros
   /system clock set time-zone-name=Europe/Stockholm
   ```

### Step 2: Configure VLAN for Home Network

```routeros
/interface vlan
add interface=ether1 name=ha vlan-id=308
```

This creates a VLAN interface for communication with your home network.

### Step 3: Configure WiFi to Connect to PHEV

```routeros
# Create security profile
/interface wireless security-profiles
add authentication-types=wpa-psk,wpa2-psk \
    management-protection=allowed \
    mode=dynamic-keys \
    name=Outlander \
    wpa2-pre-shared-key="<YOUR_PHEV_WIFI_PASSWORD>"

# Configure WiFi client
/interface wireless
set [ find default-name=wlan1 ] \
    band=2ghz-b \
    disabled=no \
    frequency=2422 \
    mode=station \
    security-profile=Outlander \
    ssid="REMOTE123456"  # Your PHEV SSID
```

**Important:** Replace `<YOUR_PHEV_WIFI_PASSWORD>` and `REMOTE123456` with your actual PHEV WiFi credentials.

### Step 4: Configure DHCP Clients

```routeros
/ip dhcp-client
# PHEV WiFi connection
add add-default-route=no \
    interface=wlan1 \
    use-peer-dns=no \
    use-peer-ntp=no \
    script=":local messagetrue \
        \"{\\\"wifiphevbound\\\":\\\"true\\\"}\"\r\
    :local messagefalse \
        \"{\\\"wifiphevbound\\\":\\\"false\\\"}\"\r\
    :local broker \"homeassistantmqtt\"\r\
    :local topic \"mikrotik/phev/wifiphevbound\"\r\
    \r\
    :if (\$bound=1) do={\r\
        /log error \"WiFi bound to PHEV\"\r\
        /iot mqtt publish broker=\$broker topic=\$topic message=\$messagetrue\r\
    } else={\r\
        /log error \"WiFi disconnected from PHEV\"\r\
        /iot mqtt publish broker=\$broker topic=\$topic message=\$messagefalse\r\
    }"

# Home network connection (VLAN)
add add-default-route=no interface=ha
```

This script publishes MQTT messages when the PHEV WiFi connection is established or lost.

### Step 5: Configure MQTT Broker

```routeros
/iot mqtt brokers
add address=192.168.1.197 \
    client-id=mikrotik \
    name=homeassistantmqtt \
    username=mikrotikmqttuser \
    password="<YOUR_MQTT_PASSWORD>"
```

**Important:** Replace with your MQTT broker IP, username, and password.

### Step 6: Configure MQTT Subscriptions for WiFi Control

```routeros
/iot mqtt subscriptions
add broker=homeassistantmqtt \
    topic=homeassistant/sensor/mikrotik_sqtsqlite2garage/wifi \
    on-message=":if (\$msgData~\"\\\\{\\\"wifi\\\": \\\"disable\\\"\\\\}\") do={\
        /interface wireless disable wlan1\
    }\
    :if (\$msgData~\"\\\\{\\\"wifi\\\": \\\"enable\\\"\\\\}\") do={\
        /interface wireless enable wlan1\
    }\
    /log info \"Got data {\$msgData} from topic {\$msgTopic}\""
```

This allows remote control of the WiFi interface via MQTT messages.

### Step 7: Configure Routing

```routeros
# Route to MQTT broker via home network
/ip route
add dst-address=192.168.1.197/32 gateway=192.168.7.1 routing-table=main

# NAT configuration
/ip firewall nat
add action=masquerade chain=srcnat out-interface=wlan1
add action=masquerade chain=srcnat out-interface=ha
add action=accept chain=srcnat src-address=192.168.8.0/24
add action=accept chain=srcnat dst-address=192.168.8.0/24
```

### Step 8: (Optional) Configure Ping Monitoring Script

This script can be disabled by default but provides automatic WiFi restart on connection loss.

```routeros
/system script
add dont-require-permissions=no \
    name=carConnectionCheck \
    owner=admin \
    policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    source="/log info \"Checking PHEV connection\"\
    \n:local message \"{\\\"wifiresetscriptrun\\\":\\\"true\\\"}\"\
    \n:local broker \"homeassistantmqtt\"\
    \n:local topiccount \"phev/connection/wifiresetcount\"\
    \n:local topicwifi \"phev/connection\"\
    \n:local HOST \"192.168.8.46\"\
    \n:local PINGCOUNT 3\
    \n:local INT \"wlan1\"\
    \n\
    \n:if ([/ping address=\$HOST interface=\$INT count=\$PINGCOUNT]=0) do={\
    \n    /log error \"PHEV not reachable, restarting WiFi\"\
    \n    /iot mqtt publish broker=\$broker topic=\$topiccount message=\$message\
    \n    /iot mqtt publish broker=\$broker topic=\$topicwifi message=\"off\"\
    \n    /interface wireless disable wlan1\
    \n    :delay 60\
    \n    /interface wireless enable wlan1\
    \n    /log error \"WiFi restarted\"\
    \n} else={\
    \n    /iot mqtt publish broker=\$broker topic=\$topicwifi message=\"on\"\
    \n}"

# Scheduler (disabled by default - enable if desired)
/system scheduler
add disabled=yes \
    interval=3m \
    name=carConnectionSchedule \
    on-event="/system script run carConnectionCheck" \
    policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
```

**Note:** The ping monitoring scheduler is disabled by default. Enable it if you want automatic WiFi restarts on connection loss.

### Step 9: (Optional) Home Assistant Status Publishing

Publish MikroTik status to Home Assistant:

```routeros
/system script
add dont-require-permissions=yes \
    name=mqtt_status \
    owner=admin \
    source="#config\
    \n:local broker \"homeassistantmqtt\"\
    \n:local topic \"homeassistant/sensor/mikrotik_sqtsqlite2garage\"\
    \n:local name [/system identity get value-name=name]\
    \n:local cpuLoad [/system resource get cpu-load]\
    \n:local upTime [:tonum [/system resource get uptime]]\
    \n\
    \n/iot mqtt publish broker=\$broker \
        topic=\"\$topic\$name/cpuload\" \
        message=\"\$cpuLoad\" \
        retain=no\
    \n/iot mqtt publish broker=\$broker \
        topic=\"\$topic\$name/uptime\" \
        message=\"\$upTime\" \
        retain=no"

# Run every minute
/system scheduler
add interval=1m \
    name=Status \
    on-event="/system script run mqtt_status" \
    policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
```

---

## Configuring phev2mqtt for MikroTik

Once your MikroTik is configured, update your phev2mqtt `.env` file:

### Remote WiFi Restart

```bash
# Enable remote WiFi restart
remote_wifi_restart_enabled=true
remote_wifi_restart_topic=mikrotik/phev/restart
remote_wifi_restart_message=restart
remote_wifi_restart_min_interval=2m
```

### Remote WiFi Power Save

```bash
# Enable power save mode
remote_wifi_power_save_enabled=true
remote_wifi_control_topic=homeassistant/sensor/mikrotik_sqtsqlite2garage/wifi
remote_wifi_enable_message={"wifi": "enable"}
remote_wifi_disable_message={"wifi": "disable"}
remote_wifi_power_save_wait=5s
remote_wifi_command_wait=10s

# Set update interval > 1 minute for power save to activate
update_interval=10m
```

See [Configuration](Configuration) for complete details.

---

## Monitoring in Home Assistant

### MQTT Topics Published by MikroTik

**WiFi Connection Status:**
- Topic: `mikrotik/phev/wifiphevbound`
- Payload: `{"wifiphevbound":"true"}` or `{"wifiphevbound":"false"}`

**Connection State (if ping monitoring enabled):**
- Topic: `phev/connection`
- Payload: `on` or `off`

**WiFi Reset Count (if ping monitoring enabled):**
- Topic: `phev/connection/wifiresetcount`
- Payload: `{"wifiresetscriptrun":"true"}`

**System Status (if status script enabled):**
- Topic: `homeassistant/sensor/mikrotik_sqtsqlite2garage/<name>/cpuload`
- Topic: `homeassistant/sensor/mikrotik_sqtsqlite2garage/<name>/uptime`

### Example Home Assistant Sensors

```yaml
# configuration.yaml
mqtt:
  sensor:
    - name: "MikroTik PHEV WiFi Bound"
      state_topic: "mikrotik/phev/wifiphevbound"
      value_template: "{{ value_json.wifiphevbound }}"
      icon: mdi:wifi
    
    - name: "PHEV Connection"
      state_topic: "phev/connection"
      icon: mdi:car-connected
    
    - name: "MikroTik CPU Load"
      state_topic: "homeassistant/sensor/mikrotik_sqtsqlite2garage/SZTsqlite2garage/cpuload"
      unit_of_measurement: "%"
      icon: mdi:cpu-64-bit
```

---

## Troubleshooting

### WiFi Won't Connect to PHEV

**Check WiFi settings:**
```routeros
/interface wireless print
/interface wireless security-profiles print
```

Verify:
- SSID matches your PHEV WiFi name
- Security profile has correct password
- WiFi interface is enabled
- Frequency is correct (usually 2422 MHz)

**Check logs:**
```routeros
/log print where topics~"wireless"
```

### MQTT Not Publishing

**Check MQTT broker connection:**
```routeros
/iot mqtt brokers print
```

Should show "Connected" status.

**Manual test publish:**
```routeros
/iot mqtt publish broker=homeassistantmqtt topic=test message="hello"
```

**Check logs:**
```routeros
/log print where topics~"mqtt"
```

### Can't Reach MQTT Broker

**Check routing:**
```routeros
/ip route print
/ping 192.168.1.197
```

Verify:
- Route to MQTT broker exists
- Gateway is reachable
- VLAN interface is up

### WiFi Control Not Working

**Test MQTT subscription:**

From another device, publish test message:
```bash
mosquitto_pub -h 192.168.1.197 -u user -P pass \
  -t "homeassistant/sensor/mikrotik_sqtsqlite2garage/wifi" \
  -m '{"wifi": "disable"}'
```

**Check MikroTik logs:**
```routeros
/log print where topics~"mqtt"
```

Should see: `Got data {"wifi": "disable"} from topic ...`

---

## Best Practices

✅ **Position Bridge Near Parking Spot** - Ensure strong PHEV WiFi signal  
✅ **Use Separate VLAN** - Isolate PHEV network from main network  
✅ **Monitor Connection Status** - Set up Home Assistant alerts  
✅ **Test Power Save Mode** - Verify WiFi turns on/off correctly  
✅ **Keep RouterOS Updated** - Update to latest stable version  
✅ **Backup Configuration** - Export and save MikroTik config regularly  

---

## Next Steps

- [Configuration](Configuration) - Configure phev2mqtt WiFi management
- [Home Assistant Integration](Home-Assistant-Integration) - Monitor MikroTik status
- [Troubleshooting](Troubleshooting) - Common issues and solutions
