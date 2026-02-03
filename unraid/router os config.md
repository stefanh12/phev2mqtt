# MikroTik RouterOS Configuration for PHEV2MQTT

## MQTT Broker Configuration

```routeros
/iot mqtt brokers
add address=192.168.1.2 name=mosquitto port=1883 username=phevmqttuser password=mqttuserpassword
```

## WiFi Restart Script

This script listens to the MQTT topic for WiFi restart commands and restarts the wireless interface.

```routeros
/system script
add name=wifi-restart owner=admin policy=reboot,read,write,policy,test source=\
    "/interface wireless reset-configuration wlan1\r\
    \n:log info \"WiFi interface wlan1 restarted via MQTT command\""
```

## MQTT Topic Subscription

Subscribe to the WiFi restart topic and trigger the script:

```routeros
/iot mqtt subscribe
add broker=mosquitto topic="phev/remote_wifi_restart" qos=0 on-message=\
    "/system script run wifi-restart"
```

## Alternative: Simple Interface Reset

If you prefer a simpler approach without full reset:

```routeros
/system script
add name=wifi-restart-simple owner=admin policy=reboot,read,write source=\
    "/interface wireless disable wlan1\r\
    \n:delay 3s\r\
    \n/interface wireless enable wlan1\r\
    \n:log info \"WiFi interface wlan1 restarted via MQTT command\""
```

## Configuration Notes

- Replace `wlan1` with your actual wireless interface name
- Update the MQTT broker address, username, and password to match your setup
- The topic `phev/remote_wifi_restart` should match the `remote_wifi_restart_topic` in your `.env` file
- You can check the script execution in `/log print` on the MikroTik device

## Verification

To test the configuration:

1. Check broker connection:
   ```routeros
   /iot mqtt brokers print
   ```

2. View active subscriptions:
   ```routeros
   /iot mqtt subscribe print
   ```

3. Test the script manually:
   ```routeros
   /system script run wifi-restart
   ```

4. Monitor logs:
   ```routeros
   /log print where message~"WiFi"
   ```

