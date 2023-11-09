# phev2mqtt - Mitsubishi Outlander PHEV to MQTT gateway using a Mikrotik wifi client bridge (RBSXTsq2nD) and running on unraid

This is build on https://github.com/buxtronix/phev2mqtt with changes to be able to run on unraid on a seperate vlan 308. The docker image is available at https://hub.docker.com/r/hstefan/phev

The original code is built to run on hardware that has wifi that connects to phev. This version has home assistant, unraid and client bridge all seperate. Wifi client availability is handled by mikrotik, when then client is online (ping) mikrotik sends a mqtt connection active that's listened to. Connection check is done every 3 minutes with ping from client bridge. The RBSXTsq2nD with my 2020 phev is really stable and only goes down when the car is not in wifi range. 
Max connection time has been removed since it was meant to handle wifi issues that the RBSXTsq2nD does not have or handles by the ping script.


Tested against a MY20 vehicle





