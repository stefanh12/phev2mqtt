CA xml template for running the phev2mqtt on docker in unraid instead of on HA. 

You can install this by adding the xml file to /boot/config/plugins/dockerMan/templates-user on your Unraid USB. then in CA add container from templates. There you will find it.

The container image is located https://hub.docker.com/r/hstefan/phev/.


The setup is only tested with one setup. That is Home Assistant running on a pi4, Unraid and Mikrotik SXTsq Lite2. 

The SXTsq needs the package add-on mqtt, this enables home assistant to automate actions towards phev2mqtt depending on the connection. Most common issue is that the wifi is lost from phev and the link goes down. Restarting the connection is then usually needed from phev2mqtt also and not only on SXT.


