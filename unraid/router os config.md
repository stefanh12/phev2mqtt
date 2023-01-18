/interface vlan
add interface=ether1 name=ha vlan-id=308
/interface wireless security-profiles
set [ find default=yes ] supplicant-identity=MikroTik
add authentication-types=wpa-psk,wpa2-psk management-protection=allowed mode=\
    dynamic-keys name=Outlander supplicant-identity=MikroTik \
    wpa-pre-shared-key=PASSWORD wpa2-pre-shared-key=PASSWORD
/interface wireless
set [ find default-name=wlan1 ] band=2ghz-b disabled=no frequency=2422 \
    mac-address=1A:0D:2E:9D:D1:D0 security-profile=Outlander ssid=REMOTE_SSID
/ip neighbor discovery-settings
set discover-interface-list=!dynamic
/interface detect-internet
set detect-interface-list=all internet-interface-list=all lan-interface-list=\
    all wan-interface-list=all
/iot mqtt brokers
add address=x.x.x.x client-id=mikrotik name=user password=\
    password username=mikrotikmqttuser
/ip dhcp-client
add add-default-route=no disabled=no interface=wlan1 script=":local messagetrue \
    \\ \r\
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
add interval=3m name=carConnectionSchedule on-event=\
    "/system script run carConnectionCheck" policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon \
    start-time=startup
/system script
add dont-require-permissions=no name=carConnectionCheck owner=admin policy=\
    ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon source="/lo\
    g info \"script start\"\r\
    \n\r\
    \n:local message \\ \r\
    \n   \"{\\\"wifiresetscriptrun\\\":\\\"true\\\"}\"\r\
    \n:local broker \"homeassistantmqtt\"\r\
    \n:local topic \"mikrotik/phev/wifiresetcount\"\r\
    \n:local HOST \"192.168.8.46\"\r\
    \n:local PINGCOUNT 3\r\
    \n:local INT \"wlan1\" \r\
    \n/log error \"script init complete\"\r\
    \n\r\
    \n:if ([/ping address=\$HOST interface=\$INT count=\$PINGCOUNT]=0) do={\r\
    \n:global name=\"tunnel_car\" 0\r\
    \n/log error \"\$INT is down\"\r\
    \n/interface wireless disable wlan1\r\
    \n/log error \"MitsubihiWiFI DISABLED\"\r\
    \n:delay 2\r\
    \n/interface wireless enable wlan1\r\
    \n:delay 4\r\
    \n/iot mqtt publish broker=\$broker topic=\$topic message=\$message\r\
    \n /log error \"MitsubihiWiFI ENABLED\"\r\
    \n} else={\r\
    \n:global name=\"tunnel_car\" 1\r\
    \n}\r\
    \n\r\
    \n"
