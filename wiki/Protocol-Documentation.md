# Protocol Documentation

This document describes the Mitsubishi Outlander PHEV remote communication protocol (MY18 and newer).

> **Note:** This information was determined through analysis, reverse engineering, and by examining the [phev-remote](https://github.com/phev-remote) code.

## Network Transport

### Connection Details

- **SSID Format**: `REMOTExxxxxx` (where `xxxxxx` is specific to your vehicle)
- **Vehicle IP**: `192.168.8.46` (fixed)
- **Client IP**: `192.168.8.47` (usually assigned via DHCP)
- **Protocol**: TCP on port `8080`
- **Data Format**: Binary packets

### Connection Sequence

1. Client connects to PHEV WiFi network
2. Client receives IP via DHCP (usually 192.168.8.47)
3. Client opens TCP connection to 192.168.8.46:8080
4. Binary packet exchange begins

---

## Packet Format

All packets follow this structure:

```text
|  1 byte   | 1 byte   | 1 byte | n bytes ... |   1 byte     |
|   [Type]    [Length]    [Ack]       [Data]     [Checksum]  |
```

### Type (1 byte)

The packet type identifier. Most types have corresponding response types.
See [Packet Types](#packet-types) below.

### Length (1 byte)

The length of the packet, **minus two**.

**Example:** A packet containing 2 data bytes has a length field of `4`.

### Ack (1 byte)

The acknowledgement field:
- `0x00` - Request packet
- `0x01` - Acknowledgement packet

### Data (Variable length)

The payload for the packet. Content depends on the packet type.

**Length calculation:** `Length field - 4`

### Checksum (1 byte)

Basic packet integrity check. Calculated by:
1. Sum all previous bytes in the packet
2. Take the least significant octet

**Example:**
```
Packet: f3 04 00 20
Sum:    0xf3 + 0x04 + 0x00 + 0x20 = 0x117
Checksum: 0x17 (least significant octet)
Final packet: f3 04 00 20 17
```

---

## Security and XOR Encoding

Most packets are obfuscated using XOR encoding with a rolling code scheme.

### Algorithm Overview

Reverse engineered by @zivillian ([GitHub issue](https://github.com/buxtronix/phev2mqtt/issues/16)).

**Steps:**

1. **Initialization** - Car sends init packet (0x5e or 0x4e)
2. **Key Derivation** - Derive `security_key` from init packet payload
3. **Key Map Generation** - Generate 256-byte key map from security key
4. **Index Tracking** - Maintain separate indices for send/receive (`s_num`, `r_num`)
5. **Packet Encoding** - XOR each byte with current key map value

### Security Initialization

Shortly after connection, the car sends an initialization packet:
- **MY18+**: Type `0x5e`
- **MY14**: Type `0x4e`

The 12-byte payload is used to derive the `security_key` and generate the `key_map` containing 256 distinct byte values.

**Indices:**
- `s_num` - Send index (client to car)
- `r_num` - Receive index (car to client)

Both start at zero.

### Packet Encoding/Decoding

After security initialization, all packets are XORed:
- **Sending to car**: XOR with `key_map[s_num]` (called `s_key`)
- **Receiving from car**: XOR with `key_map[r_num]` (called `r_key`)

XOR is symmetrical, so encoding and decoding use the same operation.

### Index Incrementing

Indices increment on specific packet types:

**Receive index (`r_num`):**
- Increments when car sends packet type `0x6f`
- New `r_key` becomes `key_map[r_num]`
- Applied to packets **after** the `0x6f` packet

**Send index (`s_num`):**
- Increments when client sends packet type `0xf6`
- New `s_key` becomes `key_map[s_num]`
- Applied to packets **after** the `0xf6` packet

**Rollover**: Both indices wrap to zero after `0xff`.

---

## Packet Types

### Summary Table

| Value | Name              | Direction     | Description                          |
|-------|-------------------|---------------|--------------------------------------|
| 0xf3  | Ping Request      | client → car  | Ping/Keepalive request               |
| 0x3f  | Ping Response     | car → client  | Ping/Keepalive response              |
| 0x6f  | Register Changed  | car → client  | Notify register has been updated     |
| 0xf6  | Register Update   | client → car  | Client register change ack/set       |
| 0x5e  | Security Init     | car → client  | Initialize security keys (MY18+)     |
| 0xe5  | Security Ack      | client → car  | Ack security keys (MY18+)            |
| 0x4e  | Security Init     | car → client  | Initialize security keys (MY14)      |
| 0xe4  | Security Ack      | client → car  | Ack security keys (MY14)             |
| 0x2f  | Keepalive Request | car → client  | Car checks client presence           |
| 0xbb  | Bad XOR           | car → client  | XOR value incorrect                  |
| 0xcc  | Bad XOR           | car → client  | XOR value incorrect                  |

### 0xf3 - Ping Request

**Format:**
```text
[f3][04][00][<seq>][00][<cksum>]
```

**Description:**
- Keepalive sent from client to car
- Initial XOR is `0x00` until `0x5e` packet received
- `<seq>` increments from `0x00` to `0x63`, then wraps to `0x00`

### 0x3f - Ping Response

**Format:**
```text
[3f][04][01][<seq>][00][<cksum>]
```

**Description:**
- Response to `0xf3` packet
- `<seq>` matches the request
- Car chooses XOR value for future register updates

### 0x6f - Register Changed

**Format:**
```text
[6f][<len>][<ack>][<register>][<data>][<cksum>]
```

**Description:**
- Notifies client that a register has changed
- `<ack>`:
  - `0x00` - Notification of register value
  - `0x01` - Response to client register change request
- `<register>` - One-byte register number
- `<data>` - Variable length, register-specific

See [Read Registers](#read-registers) for details.

### 0xf6 - Register Update/Ack

**Format:**
```text
[f6][<len>][<ack>][<register>][<data>][<cksum>]
```

**Description:**
- Client notifies car of register change or acknowledgement
- `<ack>`:
  - `0x00` - Setting new register value
  - `0x01` - Acknowledging received register update (data = `0x00`)

See [Write Registers](#write-registers) for details.

### 0x5e - Init Security (MY18+)

**Format:**
```text
[5e][0c][00][<data>][<cksum>]
```

**Description:**
- Sent by car after initial ping/keepalive exchanges
- `<data>` is 12 bytes for XOR key map initialization
- **Not** XOR encrypted

### 0xe5 - Init Ack (MY18+)

**Format:**
```text
[e5][04][01][0100][<cksum>]
```

**Description:**
- Response to `0x5e` packet
- Always the same content
- **Not** XOR encrypted

### 0xbb - Bad XOR

**Format:**
```text
[bb][06][01][<unknown>][<exp>]
```

**Description:**
- Sent when received packet has incorrect XOR
- `<exp>` field contains expected XOR value
- Can re-send packet with corrected XOR as workaround

---

## Read Registers

Registers contain vehicle state information. Sent from car to client via `0x6f` packets.

### Register Layout

There are two register layout types (referred to as A/B).

### Key Registers

| Register | Name                    | Length | Description                              |
|----------|-------------------------|--------|------------------------------------------|
| 0x02     | Battery Warning         | 4      | Battery warning status                   |
| 0x04     | Charge Timer Settings   | 20/1   | Charge timer configuration               |
| 0x05     | Climate Timer Settings  | 16/1   | Climate timer configuration              |
| 0x0b     | Parking Light Status    | 1      | Parking lights on/off                    |
| 0x10     | Pre-conditioning Status | 3/1    | Heater/AC/windscreen state               |
| 0x12     | Time Sync               | 7      | Current car time                         |
| 0x15     | VIN / Registration      | 20     | VIN and number of registered clients     |
| 0x17     | Charge Timer State      | 1      | Charge timer enabled/disabled            |
| 0x1a     | Ignition Status         | 5/2    | Ignition state (off/acc/on)              |
| 0x1c     | Aircon Mode             | 1      | Cooling/heating/windscreen               |
| 0x1d     | Battery Level / Lights  | 4      | Drive battery % and parking light status |
| 0x1e     | Charge Plug Status      | 2      | Plugged in / unplugged                   |
| 0x1f     | Charge State            | 3      | Charging status and time remaining       |
| 0x23     | Interior/Hazard Lights  | 5      | Interior and hazard light states         |
| 0x24     | Door Lock Status        | 10     | Lock, door, bonnet, boot states          |
| 0xc0     | ECU Version             | 13     | ECU software version string              |

### Register Details

#### 0x05 - Climate Timer Settings

**Length:** 16 bytes

| Byte(s) | Description |
|---------|-------------|
| 0       | Unknown     |
| 1-3     | Timer 1     |
| 4-6     | Timer 2     |
| 7-9     | Timer 3     |
| 10-12   | Timer 4     |
| 13-15   | Timer 5     |

**Timer Encoding** (3 bytes each, big-endian to little-endian):

Example: `0xa5f0afe` → `0xfef0a5`

```text
0xfef0a5 = 11111110 11110000 10100101
```

| Bits  | Description                                      |
|-------|--------------------------------------------------|
| 0-1   | Duration: 0=10min, 1=20min, 2=30min, ?3=40min?   |
| 2     | Sunday                                           |
| 3     | Monday                                           |
| 4     | Tuesday                                          |
| 5     | Wednesday                                        |
| 6     | Thursday                                         |
| 7     | Friday                                           |
| 8     | Saturday                                         |
| 9-11  | Minute: 0=0, 1=10, 2=20, 3=30, 4=40, 5=50        |
| 12-16 | Hour: 0=0, 1=1, ... 23=23                        |
| 17    | Enabled                                          |
| 18    | Disabled                                         |
| 19-23 | Unknown (unused?)                                |

#### 0x10 - Pre-conditioning Status

**Length:** 3 bytes

| Byte | Description                                                     |
|------|-----------------------------------------------------------------|
| 0    | State: 0=off, 2=on, 3=cancelled (door open/low battery)        |
| 1    | Unknown                                                         |
| 2    | Unknown                                                         |

**Examples:**
- `02b00b` - Windscreen on/10min
- `030000` - Cancelled due to door opening
- `000000` - Not active/terminated normally

#### 0x15 - VIN / Registration State

**Length:** 20 bytes

| Byte(s) | Description                   |
|---------|-------------------------------|
| 0       | Unknown                       |
| 1-18    | VIN (ASCII string)            |
| 19      | Number of registered clients  |

#### 0x1a - Ignition Status

**Length:** 5/2 bytes

| Byte | Description                                |
|------|--------------------------------------------|
| 0    | State: 0x0=off, 0x3=acc, 0x4=on            |
| 1-4  | Unknown (currently only seen as 0x0)       |

#### 0x1c - Aircon Mode

**Length:** 1 byte

| Value | Description |
|-------|-------------|
| 0     | Unknown     |
| 1     | Cooling     |
| 2     | Heating     |
| 3     | Windscreen  |

#### 0x1d - Battery Level / Parking Lights

**Length:** 4 bytes

Example: `10000003`

| Byte | Description              |
|------|--------------------------|
| 0    | Drive battery level (%)  |
| 1    | Unknown                  |
| 2    | Parking light status 0/1 |
| 3    | Unknown                  |

#### 0x1e - Charge Plug Status

**Length:** 2 bytes

| Value  | Description                |
|--------|----------------------------|
| 0x0000 | Unplugged                  |
| 0x0001 | Plugged in, not charging   |
| 0x0002 | Plugged in                 |
| 0x0202 | Charging                   |

#### 0x1f - Charging Status

**Length:** 3 bytes

| Byte(s) | Description                              |
|---------|------------------------------------------|
| 0       | Charge status: 0=not charging, 1=charging|
| 1-2     | Charge time remaining (minutes)          |

#### 0x24 - Door / Lock Status

**Length:** 10 bytes

| Byte | Description                              |
|------|------------------------------------------|
| 0    | Lock status: 1=locked, 2=unlocked        |
| 1    | Unknown                                  |
| 2    | Unknown                                  |
| 3    | Driver door (1=open, 0=closed)           |
| 4    | Front passenger door (1=open, 0=closed)  |
| 5    | Rear right door (1=open, 0=closed)       |
| 6    | Rear left door (1=open, 0=closed)        |
| 7    | Boot/trunk (1=open, 0=closed)            |
| 8    | Bonnet/hood (1=open, 0=closed)           |
| 9    | Headlight state                          |

---

## Write Registers

Commands sent from client to car via `0xf6` packets.

### Command Register Table

| Register | Name                       | Values/Description              |
|----------|----------------------------|---------------------------------|
| 0x05     | Sync Time                  | 8 bytes (see below)             |
| 0x06     | Request Updated State      | 0x03                            |
| 0x0a     | Set Head Lights            | 0x01=on, 0x02=off               |
| 0x0b     | Set Parking Lights         | 0x01=on, 0x02=off               |
| 0x0e     | Save Settings              | Sent after 0x0f command         |
| 0x0f     | Update Settings            | Various (see notes)             |
| 0x10     | Register WiFi Client       | 0x01                            |
| 0x13     | Reset PreAC State          | 0x01                            |
| 0x15     | Unregister WiFi Client     | 0x01                            |
| 0x17     | Cancel Charge Timer        | 0x01                            |
| 0x19     | Set Charge Timer Schedule  | (complex, see notes)            |
| 0x1a     | Set Climate Timer Schedule | 16 bytes (see below)            |
| 0x1b     | Set Climate State          | 4 bytes (see below)             |

### Command Details

#### 0x05 - Sync Time

**Length:** 8 bytes

| Byte | Description             |
|------|-------------------------|
| 0    | Year - 2000             |
| 1    | Month (1-12)            |
| 2    | Day of month (1-31)     |
| 3    | Hour (0-23)             |
| 4    | Minute (0-59)           |
| 5    | Second (0-59)           |
| 6    | Day of week (0-6)       |
| 7    | Device rooted (0 or 1)  |

#### 0x10 - Register WiFi Client

**Usage:** First-time registration

Sent when car is in registration mode. Car registers the client's MAC address.

#### 0x15 - Unregister WiFi Client

**Usage:** Remove registration

Removes WiFi registration for the client based on MAC address.

#### 0x1a - Set Climate Timer

**Length:** 16 bytes

| Byte(s) | Description |
|---------|-------------|
| 0-2     | Timer 1     |
| 3-5     | Timer 2     |
| 6-8     | Timer 3     |
| 9-11    | Timer 4     |
| 12-14   | Timer 5     |
| 15      | Unknown     |

Timer encoding same as register 0x05 (see above).

#### 0x1b - Set Climate State

**Format:**
```text
[02][state][duration][start]
```

| Byte | Description                                      |
|------|--------------------------------------------------|
| 0    | Always 0x02                                      |
| 1    | State: 0x01=cooling, 0x02=heating, 0x03=windscreen |
| 2    | Duration: 0x00=10min, 0x01=20min, 0x02=30min     |
| 3    | Delay: 0x00=now, 0x01=5min, 0x02=10min           |

**Examples:**
- Windscreen for 10 mins now: `02030000`
- Heat in 5 mins for 10 mins: `02020001`
- Cool for 20 mins now: `02010100`

---

## Development Tools

### Hex Decoder

```bash
./phev2mqtt decode <hex_packet>
```

Decodes and displays packet structure.

### Packet Capture

```bash
./phev2mqtt pcap <capture_file>
```

Analyzes PCAP files containing PHEV protocol traffic.

---

## References

- [phev-remote](https://github.com/phev-remote) - Original protocol research
- [buxtronix/phev2mqtt](https://github.com/buxtronix/phev2mqtt) - Original implementation
- [XOR algorithm issue](https://github.com/buxtronix/phev2mqtt/issues/16) - XOR encoding details

---

## Next Steps

- [Development](Development) - Build and contribute to phev2mqtt
- [Configuration](Configuration) - Configure phev2mqtt settings
- [Troubleshooting](Troubleshooting) - Debug protocol issues
