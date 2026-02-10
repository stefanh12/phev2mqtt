# Home Assistant Integration

phev2mqtt provides seamless integration with Home Assistant through MQTT auto-discovery. All PHEV sensors, switches, and controls appear automatically without manual configuration.

## Features

✅ **Automatic Discovery** - No manual configuration required  
✅ **Real-time Updates** - Live vehicle status updates  
✅ **Full Control** - Climate, lights, charging, and more  
✅ **VIN-Based Discovery** - Optional immediate discovery on startup  
✅ **Device Grouping** - All entities grouped under one device  

---

## Prerequisites

Before setting up Home Assistant integration:

1. **MQTT Broker** - Mosquitto or compatible MQTT broker
2. **Home Assistant MQTT Integration** - Configured and connected to your broker
3. **phev2mqtt** - Installed and running (see [Installation](Installation))

---

## Quick Setup

### Step 1: Install MQTT Integration in Home Assistant

If you haven't already:

1. **Navigate to** Settings → Devices & Services
2. **Click** "Add Integration"
3. **Search** for "MQTT"
4. **Configure** with your MQTT broker details:
   - Broker: `192.168.1.2` (your MQTT broker IP)
   - Port: `1883` (or your broker port)
   - Username: Your MQTT username
   - Password: Your MQTT password

### Step 2: Start phev2mqtt

With default settings, Home Assistant discovery is automatically enabled.

```bash
# In your .env file (default settings work fine)
mqtt_server=tcp://192.168.1.2:1883
mqtt_username=phevmqttuser
mqtt_password=your_password
```

### Step 3: Verify Discovery

Within a few seconds of phev2mqtt connecting to your PHEV:

1. **Navigate to** Settings → Devices & Services → MQTT
2. **Look for** "Mitsubishi Outlander PHEV" under discovered devices
3. **Click** the device to see all entities

**Discovery happens automatically** - no manual configuration needed!

---

## Optional: VIN-Based Immediate Discovery

By default, discovery happens after the first PHEV connection. To enable immediate discovery on startup:

```bash
# In your .env file
vehicle_vin=JA4J24A58KZ123456  # Your vehicle's VIN
```

**Benefits:**
- Home Assistant entities available immediately
- No delay waiting for PHEV connection
- Consistent device ID across restarts

**How to find your VIN:**
- Vehicle registration documents
- Dashboard VIN plate (driver's side)
- Check PHEV app or vehicle manual

---

## Available Entities

### Sensors

**Battery & Power**
- `sensor.phev_battery` - Battery level (%)
- `sensor.phev_battery_range` - Electric range (km/miles)
- `sensor.phev_charge_plug_state` - Charge plug status

**Doors & Locks**
- `binary_sensor.phev_locked` - Vehicle lock status
- `binary_sensor.phev_front_left_door` - Front left door
- `binary_sensor.phev_front_right_door` - Front right door
- `binary_sensor.phev_rear_left_door` - Rear left door
- `binary_sensor.phev_rear_right_door` - Rear right door
- `binary_sensor.phev_bonnet` - Hood/bonnet
- `binary_sensor.phev_boot` - Trunk/boot

**Climate**
- `sensor.phev_interior_temperature` - Cabin temperature
- `binary_sensor.phev_heater_operating` - Heater status
- `binary_sensor.phev_cooler_operating` - A/C status

**Charging**
- `binary_sensor.phev_charger_connected` - Charger connection status
- `binary_sensor.phev_charging` - Actively charging
- `sensor.phev_charge_remaining_time` - Time until full charge

**Lights**
- `light.phev_head_lights` - Headlights state
- `light.phev_park_lights` - Parking lights state

### Controls (Switches)

**Climate Control**
- `switch.phev_heat` - Start/stop heater
- `switch.phev_cool` - Start/stop air conditioning
- `switch.phev_windscreen` - Windscreen defroster

**Charging**
- `switch.phev_disable_charge_timer` - Override charge timer

**Other**
- `switch.phev_eco_mode` - ECO mode toggle

---

## Example Lovelace Card

phev2mqtt includes a complete Lovelace card example. See [lovelace.yaml](../lovelace.yaml) in the repository.

### Basic Card Configuration

```yaml
type: entities
title: Mitsubishi Outlander PHEV
entities:
  - entity: sensor.phev_battery
    name: Battery Level
  - entity: sensor.phev_battery_range
    name: Electric Range
  - entity: binary_sensor.phev_locked
    name: Locked
  - entity: binary_sensor.phev_charging
    name: Charging
  - entity: switch.phev_heat
    name: Heater
  - entity: switch.phev_cool
    name: Air Conditioning
```

### Advanced Picture Elements Card

For a visual representation with car top-down view:

```yaml
type: picture-elements
image: /local/car-top-view.png
elements:
  - type: state-icon
    entity: light.phev_head_lights
    style:
      top: 5%
      left: 50%
  - type: state-icon
    entity: binary_sensor.phev_front_left_door
    icon: mdi:car-door
    style:
      top: 50%
      left: 7%
  - type: state-icon
    entity: binary_sensor.phev_locked
    style:
      top: 60%
      left: 50%
  - type: state-icon
    entity: sensor.phev_battery
    style:
      top: 80%
      left: 90%
```

See the complete [lovelace.yaml](../lovelace.yaml) for the full interactive dashboard.

---

## MQTT Topics

Understanding the MQTT topic structure can help with debugging or custom integrations.

### Topic Structure

All topics follow the pattern:
```
<prefix>/<entity>
```

**Default prefix**: `phev`

### Examples

**Sensor Topics:**
- `phev/battery` - Battery level
- `phev/battery_range` - Electric range
- `phev/temperature` - Interior temperature
- `phev/doors/front_left` - Front left door status

**Control Topics:**
- `phev/climate/heat/set` - Set heater (publish `ON` or `OFF`)
- `phev/climate/cool/set` - Set A/C (publish `ON` or `OFF`)
- `phev/lights/head/set` - Set headlights (publish `ON` or `OFF`)

**Status Topics:**
- `phev/availability` - Connection status (`online` or `offline`)
- `phev/connected` - PHEV connection state

### Custom MQTT Topic Prefix

To change the default `phev` prefix:

```bash
# In your .env file
mqtt_topic_prefix=my_phev
```

All topics will then use `my_phev/` instead of `phev/`.

---

## Automations

### Example: Notify When Charging Completes

```yaml
automation:
  - alias: "PHEV Charge Complete"
    trigger:
      - platform: state
        entity_id: binary_sensor.phev_charging
        from: "on"
        to: "off"
    condition:
      - condition: state
        entity_id: sensor.phev_battery
        state: "100"
    action:
      - service: notify.mobile_app
        data:
          title: "PHEV Charged"
          message: "Your Outlander PHEV is fully charged!"
```

### Example: Auto-Heat Before Morning Commute

```yaml
automation:
  - alias: "PHEV Pre-Heat Morning"
    trigger:
      - platform: time
        at: "07:00:00"
    condition:
      - condition: state
        entity_id: binary_sensor.phev_charger_connected
        state: "on"
      - condition: numeric_state
        entity_id: sensor.phev_interior_temperature
        below: 15
    action:
      - service: switch.turn_on
        target:
          entity_id: switch.phev_heat
```

### Example: Alert on Door Open When Locked

```yaml
automation:
  - alias: "PHEV Door Alert"
    trigger:
      - platform: state
        entity_id:
          - binary_sensor.phev_front_left_door
          - binary_sensor.phev_front_right_door
          - binary_sensor.phev_rear_left_door
          - binary_sensor.phev_rear_right_door
        to: "on"
    condition:
      - condition: state
        entity_id: binary_sensor.phev_locked
        state: "on"
    action:
      - service: notify.mobile_app
        data:
          title: "PHEV Security Alert"
          message: "A door was opened while the vehicle is locked!"
```

### Example: Low Battery Notification

```yaml
automation:
  - alias: "PHEV Low Battery"
    trigger:
      - platform: numeric_state
        entity_id: sensor.phev_battery
        below: 20
    condition:
      - condition: state
        entity_id: binary_sensor.phev_charger_connected
        state: "off"
    action:
      - service: notify.mobile_app
        data:
          title: "PHEV Low Battery"
          message: "Battery is at {{ states('sensor.phev_battery') }}%. Consider charging."
```

---

## Troubleshooting

### Entities Not Appearing

**Check MQTT Integration:**
1. Navigate to Settings → Devices & Services → MQTT
2. Verify "Connected" status
3. Check broker IP and credentials

**Check phev2mqtt Logs:**
```bash
# Docker Compose
docker-compose logs | grep "discovery"

# Unraid
docker logs phev2mqtt | grep "discovery"
```

Look for messages like:
```
Publishing Home Assistant discovery for sensor: battery
Publishing Home Assistant discovery for switch: heat
```

**Manual Discovery Trigger:**
Restart phev2mqtt to re-publish discovery messages:
```bash
docker-compose restart
```

### Entities Show "Unavailable"

This means Home Assistant sees the entities but phev2mqtt is not connected to the PHEV.

**Check:**
1. Is phev2mqtt running? `docker ps | grep phev2mqtt`
2. Is WiFi connected to PHEV? Check logs for "Connected to PHEV"
3. Is PHEV WiFi active? (Vehicle must be within range)

**Monitor availability:**
```bash
# Subscribe to availability topic
mosquitto_sub -h 192.168.1.2 -u user -P pass -t "phev/availability"
```

Should show `online` when connected.

### Wrong Entity States

**Clear MQTT Retained Messages:**

Old retained messages can cause stale data.

```bash
# Clear all phev topics
mosquitto_pub -h 192.168.1.2 -u user -P pass -t "phev/#" -r -n
```

**Restart phev2mqtt** to republish fresh state.

### Duplicate Entities

This happens when the device identifier changes.

**Solution:**
1. Set `vehicle_vin` in `.env` for consistent device ID
2. Remove old device from Home Assistant
3. Restart phev2mqtt

---

## Advanced: Custom Sensors

Create template sensors based on PHEV data:

```yaml
# configuration.yaml
template:
  - sensor:
      - name: "PHEV Battery Charge Time Remaining Hours"
        unit_of_measurement: "h"
        state: >
          {{ (states('sensor.phev_charge_remaining_time') | int / 60) | round(1) }}
      
      - name: "PHEV Total Range"
        unit_of_measurement: "km"
        state: >
          {{ (states('sensor.phev_battery_range') | int + 50) }}
        # Assumes 50km fuel range for example
```

---

## Next Steps

- [Configuration](Configuration) - Customize phev2mqtt behavior
- [WiFi Management](WiFi-Management) - Optimize connection reliability
- [MikroTik Integration](MikroTik-Integration) - Advanced WiFi bridge setup
- [Troubleshooting](Troubleshooting) - Solve common issues
