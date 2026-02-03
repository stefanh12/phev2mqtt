#!/bin/bash

# Check for .env file in config directory and load it
if [ -f "/config/.env" ]; then
    echo "Loading configuration from /config/.env"
    # Export all variables from .env file
    set -a
    source /config/.env
    set +a
    
    # Check if the .env file has been edited by the user
    if grep -q "mqtt_password=your_secure_password_here" /config/.env; then
        echo ""
        echo "=========================================="
        echo "  ERROR: Default .env file detected"
        echo "=========================================="
        echo ""
        echo "The .env file has not been edited!"
        echo ""
        echo "Please update /config/.env with your actual:"
        echo "  - mqtt_password (REQUIRED)"
        echo "  - mqtt_server (if different)"
        echo "  - mqtt_user (if different)"
        echo "  - route_add (if needed)"
        echo ""
        echo "Container will exit now."
        echo "Restart after editing the configuration."
        echo "=========================================="
        echo ""
        exit 1
    fi
else
    # Create example .env file if it doesn't exist
    if [ -d "/config" ] && [ ! -f "/config/.env" ]; then
        echo "=========================================="
        echo "  Creating example .env file"
        echo "=========================================="
        cat > /config/.env << 'EOF'
# phev2mqtt Configuration
# Edit this file with your actual credentials and settings

# ==========================================
# MQTT Broker Configuration (REQUIRED)
# ==========================================
mqtt_server=192.168.1.2:1883
mqtt_user=phevmqttuser
mqtt_password=your_secure_password_here

# ==========================================
# PHEV Registration Mode
# ==========================================
# Set to 'true' ONLY when registering a new vehicle
# After successful registration, change back to 'false'
phev_register=false

# ==========================================
# Debug Logging
# ==========================================
# Set to 'true' to enable detailed debug logging
# This shows all connection attempts, MQTT messages, and protocol details
# Set to 'false' for normal info level logging (recommended for production)
debug=false

# ==========================================
# Network Routing
# ==========================================
# IP address of the gateway that routes to your PHEV network (192.168.8.0/24)
# This is typically your router's IP address on the main network
# Example: If your router is at 192.168.1.1, set this to 192.168.1.1
route_add=192.168.1.1

# ==========================================
# Additional Arguments (optional)
# ==========================================
# Any extra command-line arguments to pass to phev2mqtt
# Leave empty if not needed
extra_add=""
EOF
        echo ""
        echo "✓ Example .env file created at /config/.env"
        echo ""
        echo "IMPORTANT: Please follow these steps:"
        echo "  1. Edit /config/.env with your actual settings"
        echo "  2. Update mqtt_password with your MQTT password"
        echo "  3. Restart the container to apply changes"
        echo ""
        echo "Note: Container is currently using environment"
        echo "      variables from the Docker template."
        echo "=========================================="
    fi
fi

export CONNECT_DEBUG=$debug
export CONNECT_phev_register=$phev_register
export CONNECT_mqtt_server=$mqtt_server
export CONNECT_mqtt_user=$mqtt_user
export CONNECT_mqtt_password=$mqtt_password
export CONNECT_route_add=$route_add
export CONNECT_extra_add=$extra_add


echo "Using the following environment variables:"
echo "debug=$CONNECT_DEBUG"
echo "phev_register=$CONNECT_phev_register"
echo "mqtt_server=$CONNECT_mqtt_server"
#echo "mqtt_user=$CONNECT_mqtt_user"
#echo "mqtt_password=$CONNECT_mqtt_password"
echo "route_add=$CONNECT_route_add"
echo "extra_add=$CONNECT_extra_add"

# Add route if configured
if [[ -n "$CONNECT_route_add" ]]; then
    echo "Adding route to 192.168.8.0/24 via gateway $CONNECT_route_add"
    route add -net 192.168.8.0 netmask 255.255.255.0 gw $CONNECT_route_add eth0
else
    echo "Warning: route_add not set. This may cause connectivity issues with the PHEV (192.168.8.0/24)"
fi

# Set log level based on debug flag
if [[ $CONNECT_DEBUG == "true" ]]; then
    LOG_LEVEL="debug"
    echo "Debug logging enabled (verbosity=debug)"
else
    LOG_LEVEL="info"
fi

if [[ $CONNECT_phev_register == "true" ]]
then
	echo "Registering client with PHEV"
	exec /usr/src/app/phev2mqtt/phev2mqtt client register --verbosity "$LOG_LEVEL"
else
    echo "Starting phev2mqtt with log level: $LOG_LEVEL"
    exec /usr/src/app/phev2mqtt/phev2mqtt \
        client \
        mqtt \
        --verbosity "$LOG_LEVEL" \
        --mqtt_server "tcp://$CONNECT_mqtt_server/" \
        --mqtt_username "$CONNECT_mqtt_user" \
        --mqtt_password "$CONNECT_mqtt_password" \
        $CONNECT_extra_add
fi
