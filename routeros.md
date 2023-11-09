
/interface vlan
add interface=ether1 name=ha vlan-id=308
/interface wireless security-profiles
set [ find default=yes ] supplicant-identity=MikroTik
add authentication-types=wpa-psk,wpa2-psk management-protection=allowed mode=dynamic-keys name=Outlander supplicant-identity=MikroTik wpa-pre-shared-key=key wpa2-pre-shared-key=key
/interface wireless
set [ find default-name=wlan1 ] band=2ghz-b disabled=no frequency=2422 mac-address=11:0D:1E:1D:D1:D1 security-profile=Outlander ssid=PHEVSSID
/ip neighbor discovery-settings
set discover-interface-list=!dynamic
/interface detect-internet
set detect-interface-list=all internet-interface-list=all lan-interface-list=all wan-interface-list=all
/iot mqtt brokers
add address=IPMQTTSERVER client-id=mikrotik name=homeassistantmqtt password=MQTTSERVERPW username=MQTTUSER
/ip dhcp-client
add add-default-route=no disabled=no interface=wlan1 script=":local messagetrue \\ \r\
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
add add-default-route=no disabled=no interface=ha
/ip firewall nat
add action=masquerade chain=srcnat out-interface=wlan1
add action=masquerade chain=srcnat out-interface=ha
add action=accept chain=srcnat src-address=192.168.8.0/24
add action=accept chain=srcnat dst-address=192.168.8.0/24
/ip service
set api address=0.0.0.0/0
/system clock
set time-zone-name=Europe/Stockholm
/system scheduler
add interval=3m name=carConnectionSchedule on-event="/system script run carConnectionCheck" policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon start-time=startup
/system script
add dont-require-permissions=no name=carConnectionCheck owner=admin policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon source="/log info \"script start\"\
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
    \n/iot mqtt publish broker=\$broker topic=\$topicwifi message=\$messagenowifi\
    \n/interface wireless disable wlan1\
    \n/log error \"MitsubihiWiFI DISABLED\"\
    \n:delay 5\
    \n/interface wireless enable wlan1\
    \n /log error \"MitsubihiWiFI ENABLED\"\
    \n} else={\
    \n:global name=\"tunnel_car\" 1\
    \n/iot mqtt publish broker=\$broker topic=\$topicwifi message=\$messagewifi\
    \n\
    \n}"
