# 2026-02-03 15:53:46 by RouterOS 7.21.2
# software id = XXXXXXXXXXX
#
# model = RBSXTsq2nD
# serial number = XXXXXXXXX
/interface vlan
add interface=ether1 name=ha vlan-id=308
/interface lte apn
set [ find default=yes ] ip-type=ipv4 use-network-apn=no
/interface wireless security-profiles
set [ find default=yes ] supplicant-identity=MikroTik
add authentication-types=wpa-psk,wpa2-psk management-protection=allowed mode=\
    dynamic-keys name=Outlander supplicant-identity=MikroTik
/interface wireless
set [ find default-name=wlan1 ] band=2ghz-b disabled=no frequency=2422 \
    mac-address=XX:XX:XX:XX:XX:XX security-profile=Outlander ssid=XXXXXXXX
/iot lora servers
add address=eu.mikrotik.thethings.industries name=TTN-EU protocol=UDP
add address=us.mikrotik.thethings.industries name=TTN-US protocol=UDP
add address=eu1.cloud.thethings.industries name="TTS Cloud (eu1)" protocol=UDP
add address=nam1.cloud.thethings.industries name="TTS Cloud (nam1)" protocol=\
    UDP
add address=au1.cloud.thethings.industries name="TTS Cloud (au1)" protocol=UDP
add address=eu1.cloud.thethings.network name="TTN V3 (eu1)" protocol=UDP
add address=nam1.cloud.thethings.network name="TTN V3 (nam1)" protocol=UDP
add address=au1.cloud.thethings.network name="TTN V3 (au1)" protocol=UDP
/iot mqtt brokers
add address=192.168.1.197 client-id=mikrotik name=homeassistantmqtt username=\
    mikrotikmqttuser
/ip smb users
set [ find default=yes ] disabled=yes
/routing bgp template
set default disabled=no output.network=bgp-networks
/routing ospf instance
add disabled=no name=default-v2
/routing ospf area
add disabled=yes instance=default-v2 name=backbone-v2
/ip firewall connection tracking
set udp-timeout=10s
/ip neighbor discovery-settings
set discover-interface-list=!dynamic
/ip settings
set max-neighbor-entries=8192
/ipv6 settings
set disable-ipv6=yes max-neighbor-entries=8192
/interface detect-internet
set detect-interface-list=all internet-interface-list=all lan-interface-list=\
    all wan-interface-list=all
/interface ovpn-server server
add auth=sha1,md5 mac-address=XX:XX:XX:XX:XX:XX name=ovpn-server1
/iot mqtt subscriptions
add broker=homeassistantmqtt on-message=":if (\$msgData~\"\\\\{\\\"wifi\\\": \\\
    \"disable\\\"\\\\}\") do={/interface wireless disable wlan1}\
    \n:if (\$msgData~\"\\\\{\\\"wifi\\\": \\\"enable\\\"\\\\}\") do={/interface\
    \_wireless enable wlan1}\
    \nlog info \"Got data {\$msgData} from topic {\$msgTopic}\"" topic=\
    homeassistant/sensor/mikrotik_sqtsqlite2garage/wifi
/ip dhcp-client
# Interface not active
add add-default-route=no interface=wlan1 script=":local messagetrue \\ \r\
    \n   \"{\\\"wifiphevbound\\\":\\\"true\\\"}\"\r\
    \n:local messagefalse \\ \r\
    \n   \"{\\\"wifiphevbound\\\":\\\"false\\\"}\"   \r\
    \n:local broker \"homeassistantmqtt\"\r\
    \n:local topic \"mikrotik/phev/wifiphevbound\"\r\
    \n\r\
    \n:if (\\\$bound=1) do={\r\
    \n/log error \"script bound = 1\"\r\
    \n/iot mqtt publish broker=\$broker topic=\$topic message=\$messagetrue\r\
    \n\r\
    \n} else={\r\
    \n/log error \"script bound = 0\"\r\
    \n/iot mqtt publish broker=\$broker topic=\$topic message=\$messagefalse\r\
    \n\r\
    \n}\r\
    \n" use-peer-dns=no use-peer-ntp=no
add add-default-route=no interface=ha
/ip firewall nat
add action=masquerade chain=srcnat out-interface=wlan1
add action=masquerade chain=srcnat out-interface=ha
add action=accept chain=srcnat src-address=192.168.8.0/24
add action=accept chain=srcnat dst-address=192.168.8.0/24
/ip ipsec profile
set [ find default=yes ] dpd-interval=2m dpd-maximum-failures=5
/ip route
add disabled=no dst-address=192.168.1.197/32 gateway=192.168.7.1 \
    routing-table=main
/ip service
set api address=0.0.0.0/0
/routing bfd configuration
add disabled=no interfaces=all min-rx=200ms min-tx=200ms multiplier=5
/system clock
set time-zone-name=Europe/Stockholm
/system identity
set name=SZTsqlite2garage
/system scheduler
add disabled=yes interval=3m name=carConnectionSchedule on-event=\
    "/system script run carConnectionCheck" policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
add interval=1m name=Status on-event="/system script run  mqtt_status" policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
add name=status_config on-event="/system script run  mqtt_status_config" \
    policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
/system script
add dont-require-permissions=no name=carConnectionCheck owner=admin policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon source="/l\
    og info \"script start\"\
    \n\
    \n:local message \"{\\\"wifiresetscriptrun\\\":\\\"true\\\"}\"\
    \n:local messagenowifi \"off\"\
    \n:local messagewifi \"on\"\
    \n:local broker \"homeassistantmqtt\"\
    \n:local topiccount \"phev/connection/wifiresetcount\"\
    \n:local topicwifi \"phev/connection\"\
    \n:local HOST \"192.168.8.46\"\
    \n:local PINGCOUNT 3\
    \n:local INT \"wlan1\" \
    \n/log error \"script init complete\"\
    \n\
    \n:if ([/ping address=\$HOST interface=\$INT count=\$PINGCOUNT]=0) do={\
    \n:global name=\"tunnel_car\" 0\
    \n/log error \"\$INT is down\"\
    \n/iot mqtt publish broker=\$broker topic=\$topiccount message=\$message\
    \n/iot mqtt publish broker=\$broker topic=\$topicwifi message=\$messagenowi\
    fi\
    \n/interface wireless disable wlan1\
    \n/log error \"MitsubihiWiFI DISABLED\"\
    \n:delay 60\
    \n/interface wireless enable wlan1\
    \n /log error \"MitsubihiWiFI ENABLED\"\
    \n} else={\
    \n:global name=\"tunnel_car\" 1\
    \n/iot mqtt publish broker=\$broker topic=\$topicwifi message=\$messagewifi\
    \n\
    \n}"
add dont-require-permissions=yes name=mqtt_status_config owner=admin policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon source="#c\
    onfig\
    \n:local broker \"homeassistantmqtt\"\
    \n:local topic \"homeassistant/sensor/mikrotik_sqtsqlite2garage\"\
    \n\
    \n#gather info\
    \n:local model [/system routerboard get value-name=model]\
    \n:local name [/system identity get value-name=name]\
    \n:local board [/system resource get value-name=board-name]\
    \n:local version [/system resource get value-name=version]\
    \n:local serial [/system routerboard get serial-number]\
    \n\
    \n#compose message in JSON - CPU load\
    \n:local message  \"{\\\"device\\\": {\\\"mf\\\": \\\"MikroTik\\\", \\\"ide\
    ntifiers\\\": \\\"MT-\$serial\\\", \\\"model\\\": \\\"\$model\\\", \\\"sw\\\
    \": \\\"\$version\\\", \\\"name\\\": \\\"\$name\\\"}, \\\"unique_id\\\": \\\
    \"MT-\$serial-cpuload\\\", \\\"state_topic\\\": \\\"\$topic\$name/cpuload\\\
    \", \\\"unit_of_measurement\\\": \\\"%\\\", \\\"state_class\\\": \\\"measur\
    ement\\\", \\\"name\\\": \\\"CPU load\\\"}\"\
    \n:put \$message\
    \n\
    \n/iot mqtt publish broker=\$broker topic=\"\$topic\$name/cpuload/config\" \
    message=\$message retain=yes\
    \n/log/info message=\$message\
    \n\
    \n#Up time - \\\"board\\\": \\\"\$board\\\", \
    \n:local message  \"{\\\"device\\\": {\\\"mf\\\": \\\"MikroTik\\\", \\\"ide\
    ntifiers\\\": \\\"MT-\$serial\\\", \\\"model\\\": \\\"\$model\\\", \\\"sw\\\
    \": \\\"\$version\\\", \\\"name\\\": \\\"\$name\\\"}, \\\"unique_id\\\": \\\
    \"MT-\$serial-uptime\\\", \\\"state_topic\\\": \\\"\$topic\$name/uptime\\\"\
    , \\\"unit_of_measurement\\\": \\\"s\\\", \\\"state_class\\\": \\\"total\\\
    \", \\\"name\\\": \\\"Up time\\\"}\"\
    \n:put \$message\
    \n\
    \n/iot mqtt publish broker=\$broker topic=\"\$topic\$name/uptime/config\" m\
    essage=\$message retain=yes\
    \n/log/info message=\$message"
add dont-require-permissions=yes name=mqtt_status owner=admin policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon source="#c\
    onfig\
    \n:local broker \"homeassistantmqtt\"\
    \n:local topic \"homeassistant/sensor/mikrotik_sqtsqlite2garage\"\
    \n\
    \n#gather info\
    \n:local name [/system identity get value-name=name]\
    \n\
    \n:local cpuLoad [/system resource get cpu-load]\
    \n:local upTime [:tonum [/system resource get uptime]]\
    \n\
    \n#CPU load\
    \n:local message  \"\$cpuLoad\"\
    \n:put \$message\
    \n\
    \n/iot mqtt publish broker=\$broker topic=\"\$topic\$name/cpuload\" message\
    =\$message retain=no\
    \n/log/info message=\$message\
    \n\
    \n#Up time\
    \n:local message2 \"\$upTime\"\
    \n:put \$message2\
    \n\
    \n/iot mqtt publish broker=\$broker topic=\"\$topic\$name/uptime\" message=\
    \$message2 retain=no\
    \n/log/info message=\$message2"