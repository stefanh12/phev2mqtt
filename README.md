# phev2mqtt - Mitsubishi Outlander PHEV to MQTT gateway using a Mikrotik wifi client bridge (RBSXTsq2nD) and running on unraid

This is build on https://github.com/buxtronix/phev2mqtt and https://github.com/CodeCutterUK/phev2mqtt with changes to be able to run on unraid on a seperate vlan 308. The docker image is available at https://hub.docker.com/r/hstefan/phev

The original code is built to run on hardware that has wifi that connects to phev. This version has home assistant, unraid and client bridge all seperate. Wifi client availability is handled by mikrotik, when then client is online (ping) mikrotik sends a mqtt connection active that's listened to. Connection check is done every 3 minutes with ping from client bridge. The RBSXTsq2nD with my 2020 phev is really stable and only goes down when the car is not in wifi range. 
Max connection time has been removed since it was meant to handle wifi issues that the RBSXTsq2nD does not have or handles by the ping script.

## WiFi Restart Features

The application supports two types of WiFi restart mechanisms to handle connection issues:

### Local WiFi Restart
Automatically restarts the local WiFi interface when connection to the PHEV is lost. This is useful if phev2mqtt is running on hardware with its own WiFi interface.
- Configure via `LOCAL_WIFI_RESTART_ENABLED=true` in .env file
- Set `WIFI_RESTART_TIME` to define how long to wait before restarting (e.g., 10m)
- Note: Only works on some hardware configurations

### Remote WiFi Restart (MikroTik Integration)
Sends MQTT commands to remotely restart WiFi on external devices like MikroTik access points when connection is lost.
- Configure via `REMOTE_WIFI_RESTART_ENABLED=true` in .env file
- Set `REMOTE_WIFI_RESTART_TOPIC` to the MQTT topic your MikroTik subscribes to
- The MikroTik script will receive the restart command and reset the WiFi interface
- Useful when using a dedicated WiFi bridge like RBSXTsq2nD

See the [routeros.md](routeros.md) file for MikroTik configuration examples.

## Hot Reload Configuration

The application automatically watches the `.env` file for changes and reloads configuration without requiring a container restart. Changes are detected every 5 seconds.

**Hot-reloadable settings:**
- `update_interval` - Update request frequency
- `retry_interval` - Connection retry frequency
- `remote_wifi_restart_topic` - Remote restart MQTT topic
- `remote_wifi_restart_message` - Remote restart message

**Settings requiring restart:**
- MQTT connection settings (server, username, password)
- Debug logging level
- PHEV registration mode
- Network routing configuration

To apply changes, simply edit `/config/.env` and save. The changes will be applied within 5 seconds. Check the logs to confirm reload:
```
Configuration file changed, reloading...
Configuration reloaded successfully
```

Tested against a MY20 vehicle

Routeros config 
https://github.com/stefanh12/phev2mqtt/blob/master/routeros.md 

unraid
https://github.com/stefanh12/phev2mqtt/blob/master/unraid/XMLFile1.xml


