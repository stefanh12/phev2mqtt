#!/bin/bash
export CONNECT_DEBUG=$debug
export CONNECT_phev_register=$phev_register
export CONNECT_mqtt_server=$mqtt_server
export CONNECT_mqtt_user=$mqtt_user
export CONNECT_mqtt_password=$mqtt_password
echo "Using the following environment variables:"
echo "debug=$CONNECT_DEBUG"
echo "phev_register=$CONNECT_phev_register"
echo "mqtt_server=$CONNECT_mqtt_server"
echo "mqtt_user=$CONNECT_mqtt_user"
echo "mqtt_password=$CONNECT_mqtt_password"

if [[ "x$CONNECT_DEBUG" = "x" ]]; then
		echo "The debug variable is not set, should be set to true to sleep. Can be used to register the client with /usr/src/app/phev2mqtt/phev2mqtt client register"
fi
if [[ "x$CONNECT_phev_register" = "x" ]]; then
		echo "The phev_register variable shall be set to true to register the client with the phev"
		exit 1
fi
if [[ "x$CONNECT_mqtt_server" = "x" ]]; then
		echo "The mqtt_server variable must be set."
		exit 1
fi
if [[ "x$CONNECT_mqtt_user" = "x" ]]; then
		echo "The mqtt_user variable must be set."
		exit 1
fi
if [[ "x$CONNECT_mqtt_password" = "x" ]]; then
		echo "The CONNECT_mqtt_password variable must be set."
		exit 1
fi

if [[ $CONNECT_DEBUG == "true" ]]
then
	echo Debug mode on - sleeping indefinitely
	sleep inf
fi
if [[ $CONNECT_phev_register == "true" ]]
then
	echo register client
	/usr/src/app/phev2mqtt/phev2mqtt client register
else
    echo Starting phev2mqtt
    /usr/src/app/phev2mqtt/phev2mqtt \
        client \
        mqtt \
        --mqtt_server "tcp://$CONNECT_mqtt_server/" \
        --mqtt_username "$CONNECT_mqtt_user" \
        --mqtt_password "$CONNECT_mqtt_password"
fi