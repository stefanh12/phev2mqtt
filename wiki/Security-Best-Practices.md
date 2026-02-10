# Security Best Practices

This guide covers security considerations and best practices for deploying phev2mqtt.

## Security Status

phev2mqtt has undergone security auditing (February 2026). **Critical and high-priority vulnerabilities have been fixed.**

✅ **Command injection prevention** - WiFi restart commands validated  
✅ **MQTT password requirements** - Minimum 8 characters enforced  
✅ **MQTT server validation** - URL format and protocol verified  
✅ **MQTT topic validation** - Prevents topic injection attacks  

For complete audit details, see [SECURITY_AUDIT.md](../SECURITY_AUDIT.md).

---

## Password Security

### MQTT Password Requirements

**Hard Requirements (Enforced):**
- Minimum 8 characters
- Not a common weak password

**Recommendations:**
- Use 12+ characters
- Include mix of uppercase, lowercase, numbers, special characters
- Use a password manager
- Don't reuse passwords from other services

**Weak Passwords Blocked:**
- password
- admin123
- 12345678
- mqtt1234
- And other common patterns

### Storing Passwords

**File-Based Configuration (.env):**

✅ **Pros:**
- Not visible in Docker UI
- Easy to backup separately
- Can be encrypted at rest

⚠️ **Cons:**
- Still stored in plaintext in file
- File permissions critical

**Best practices:**
```bash
# Restrict permissions
chmod 600 /path/to/.env
chown user:user /path/to/.env

# For Unraid
chmod 600 /mnt/user/appdata/phev2mqtt/.env
```

**Environment Variable Configuration:**

⚠️ **Less secure:**
- Visible in Docker inspect
- Visible in Unraid UI (though masked)
- May be logged

**Recommendation:** Use `.env` file method for better security.

### Password Manager Integration

Consider using:
- **1Password** - supports Docker secrets
- **Bitwarden** - open source, self-hostable
- **KeePass** - local, encrypted database

---

## MQTT Broker Security

### Use TLS/SSL

**Always use encrypted connections in production:**

```bash
# In .env file
mqtt_server=ssl://mqtt.example.com:8883
# or
mqtt_server=tls://mqtt.example.com:8883
```

**Unencrypted connections (tcp://, ws://) trigger security warnings.**

### Certificate Validation

For TLS connections:
- Use valid SSL certificates (Let's Encrypt)
- Don't disable certificate verification
- Ensure hostname matches certificate

### MQTT Broker Hardening

**Mosquitto example configuration:**

```conf
# /etc/mosquitto/mosquitto.conf

# Require authentication
allow_anonymous false
password_file /etc/mosquitto/passwd

# Use TLS
listener 8883
cafile /etc/mosquitto/ca_certificates/ca.crt
certfile /etc/mosquitto/certs/server.crt
keyfile /etc/mosquitto/certs/server.key

# TLS version
tls_version tlsv1.2

# Require client certificates (optional, for mutual TLS)
require_certificate true
```

**Create password file:**
```bash
mosquitto_passwd -c /etc/mosquitto/passwd phevmqttuser
```

### MQTT ACLs (Access Control Lists)

Limit what phev2mqtt can publish/subscribe to:

```conf
# /etc/mosquitto/acl
user phevmqttuser
topic write phev/#
topic readwrite homeassistant/#
topic read mikrotik/#
```

---

## Network Security

### VLAN Isolation

**Recommended network architecture:**

1. **Management VLAN** - Home Assistant, MQTT broker
2. **IoT VLAN** - phev2mqtt container
3. **PHEV VLAN** - WiFi bridge to vehicle

**Benefits:**
- Isolate vehicle network from main network
- Control traffic flow with firewall rules
- Limit blast radius if compromised

### Firewall Rules

**Example using iptables:**

```bash
# Allow phev2mqtt to MQTT broker only
iptables -A FORWARD -s 192.168.10.0/24 -d 192.168.1.197 -p tcp --dport 1883 -j ACCEPT

# Allow phev2mqtt to PHEV network
iptables -A FORWARD -s 192.168.10.0/24 -d 192.168.8.0/24 -j ACCEPT

# Allow Home Assistant to MQTT
iptables -A FORWARD -s 192.168.1.0/24 -d 192.168.1.197 -p tcp --dport 1883 -j ACCEPT

# Drop everything else
iptables -A FORWARD -j DROP
```

### Disable Unnecessary Services

On MQTT broker:
```bash
# Disable HTTP listener if using TLS
# In mosquitto.conf, don't enable listener on port 1883

# Only enable what you need
listener 8883
protocol mqtt
```

---

## Container Security

### Run as Non-Root User

phev2mqtt requires `NET_ADMIN` capability for routing but can limit privileges:

```yaml
# docker-compose.yml
services:
  phev2mqtt:
    image: ghcr.io/stefanh12/phev2mqtt:latest
    cap_add:
      - NET_ADMIN
    cap_drop:
      - ALL
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
```

### Keep Images Updated

```bash
# Update regularly
docker-compose pull
docker-compose up -d

# Check for updates
docker images | grep phev2mqtt
```

**Subscribe to:**
- GitHub releases for notifications
- Security advisories

### Scan for Vulnerabilities

```bash
# Using Docker Scout
docker scout cves ghcr.io/stefanh12/phev2mqtt:latest

# Using Trivy
trivy image ghcr.io/stefanh12/phev2mqtt:latest
```

---

## Input Validation

phev2mqtt validates all user inputs to prevent security issues.

### WiFi Restart Command Validation

**Blocked characters:**
- `;` (semicolon)
- `|` (pipe)
- `$()` (command substitution)
- Backticks
- `>`, `<` (redirects)
- Newlines

**Allowed:**
- `&&` (AND operator) - with warning

**Example safe commands:**
```bash
# Safe
wifi_restart_command=nmcli device disconnect wlan0 && nmcli device connect wlan0

# Blocked
wifi_restart_command=nmcli device disconnect wlan0; rm -rf /
```

### MQTT Topic Validation

**Blocked in topics:**
- Wildcards (`+`, `#`) in publish topics
- Null bytes
- Control characters
- Suspicious patterns (`${`, `$()`, backticks, `<script>`, `javascript:`)

**Length limits:**
- Warning at 512 characters
- Block at 65535 characters

### MQTT Server URL Validation

**Required:**
- Valid protocol prefix
- Valid port (1-65535)
- Proper URL format

**Security warnings issued for:**
- Unencrypted connections (`tcp://`, `ws://`)

---

## Credential Management

### Environment Variables

**Best practices:**

1. **Use .env file** instead of environment variables
2. **Never commit .env to git:**
   ```bash
   # Add to .gitignore
   .env
   .env.local
   *.env
   ```

3. **Use example file for documentation:**
   ```bash
   cp .env .env.example
   # Remove actual values from .env.example
   ```

4. **Restrict file access:**
   ```bash
   chmod 600 .env
   ```

### Backup Security

When backing up configuration:

1. **Encrypt backups containing credentials:**
   ```bash
   tar czf - /path/to/phev2mqtt | gpg -e -r your@email.com > backup.tar.gz.gpg
   ```

2. **Store backups securely:**
   - Encrypted cloud storage
   - Password-protected archive
   - Separate from public backups

3. **Don't backup to public locations:**
   - No public GitHub repos
   - No unencrypted cloud shares
   - No world-readable directories

---

## MikroTik Security

### RouterOS Hardening

1. **Use strong admin password:**
   ```routeros
   /user set admin password="strong_random_password"
   ```

2. **Create separate MQTT user:**
   ```routeros
   /user add name=mqtt group=write password="mqtt_password"
   ```

3. **Restrict service access:**
   ```routeros
   /ip service disable telnet,ftp,www
   /ip service set ssh port=2222
   /ip service set winbox address=192.168.1.0/24
   ```

4. **Enable firewall:**
   ```routeros
   # Block everything except what's needed
   /ip firewall filter
   add chain=input action=accept connection-state=established,related
   add chain=input action=accept src-address=192.168.1.0/24
   add chain=input action=drop
   ```

5. **Secure MQTT credentials in scripts:**
   - Scripts are visible in plain text in RouterOS
   - Consider using environment variables or separate config
   - Limit user access to script view permissions

---

## Monitoring and Logging

### Enable Appropriate Logging

```bash
# Production - info level
log_level=info

# Debugging - debug level (temporary)
log_level=debug

# Be careful with debug logs - may contain sensitive data
```

### Log Monitoring

**Monitor for suspicious activity:**
- Repeated authentication failures
- Unusual connection patterns
- Unexpected configuration reloads
- Command injection attempts

**Example log monitoring:**
```bash
# Check for errors
docker logs phev2mqtt | grep -i error

# Check for failed connections
docker logs phev2mqtt | grep -i "failed\|timeout"

# Check for validation issues
docker logs phev2mqtt | grep -i "invalid\|blocked"
```

### Log Retention

**Best practices:**
- Rotate logs regularly
- Keep 30-90 days of logs
- Archive old logs securely
- Don't log passwords (phev2mqtt doesn't)

**Docker logging config:**
```yaml
# docker-compose.yml
services:
  phev2mqtt:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## Home Assistant Security

### MQTT Integration

1. **Use dedicated MQTT user:**
   ```bash
   mosquitto_passwd -b /etc/mosquitto/passwd homeassistant ha_password
   ```

2. **Limit permissions:**
   ```conf
   # /etc/mosquitto/acl
   user homeassistant
   topic readwrite homeassistant/#
   topic read phev/#
   ```

3. **Enable authentication in Home Assistant:**
   - Configuration → People → Users
   - Enable multi-factor authentication
   - Use strong passwords

### Auto-Discovery Security

**Risks:**
- Any device publishing to `homeassistant/#` can create entities
- Rogue devices could publish false data

**Mitigations:**
- Use MQTT ACLs to limit who can publish to discovery topics
- Monitor MQTT traffic
- Review auto-discovered devices regularly

---

## Incident Response

### If Credentials Compromised

1. **Immediately change passwords:**
   ```bash
   mosquitto_passwd /etc/mosquitto/passwd phevmqttuser
   # Update .env file
   # Restart containers
   ```

2. **Review MQTT logs:**
   ```bash
   # Look for suspicious connections
   grep phevmqttuser /var/log/mosquitto/mosquitto.log
   ```

3. **Check for unauthorized access:**
   - Review Home Assistant history
   - Check for unexpected vehicle commands
   - Review MikroTik logs

4. **Update all related credentials:**
   - Home Assistant
   - MQTT broker
   - MikroTik
   - Any other integrated services

### If Container Compromised

1. **Stop container immediately:**
   ```bash
   docker stop phev2mqtt
   ```

2. **Preserve evidence:**
   ```bash
   docker logs phev2mqtt > incident_logs.txt
   docker inspect phev2mqtt > incident_inspect.txt
   ```

3. **Investigate:**
   - Review logs for unauthorized access
   - Check for modified files
   - Scan for malware

4. **Rebuild from clean image:**
   ```bash
   docker pull ghcr.io/stefanh12/phev2mqtt:latest
   docker-compose up -d
   ```

---

## Security Checklist

Use this checklist for a secure deployment:

### Required (High Priority)

- [ ] Use strong MQTT password (12+ characters)
- [ ] Use TLS for MQTT connections
- [ ] Restrict `.env` file permissions (chmod 600)
- [ ] Enable MQTT authentication
- [ ] Keep container images updated
- [ ] Use dedicated MQTT user for phev2mqtt
- [ ] Review logs regularly

### Recommended (Medium Priority)

- [ ] Implement VLAN isolation
- [ ] Configure firewall rules
- [ ] Use MQTT ACLs
- [ ] Enable Home Assistant MFA
- [ ] Configure log rotation
- [ ] Regular security audits
- [ ] Backup encrypted

### Optional (Low Priority)

- [ ] Use mutual TLS (mTLS)
- [ ] Implement IDS/IPS
- [ ] Container vulnerability scanning
- [ ] Network traffic monitoring
- [ ] Honeypot integration

---

## Reporting Security Issues

**Please report security vulnerabilities responsibly:**

1. **Do NOT open public issue**
2. **Email maintainer directly** (see GitHub profile)
3. **Provide details:**
   - Vulnerability description
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

4. **Allow time for fix:**
   - Maintainer will acknowledge within 48 hours
   - Fix will be prioritized
   - Public disclosure after patch available

---

## Additional Resources

- [OWASP IoT Top 10](https://owasp.org/www-project-internet-of-things/)
- [Docker Security Best Practices](https://docs.docker.com/engine/security/)
- [MQTT Security Fundamentals](https://www.hivemq.com/mqtt-security-fundamentals/)
- [Home Assistant Security](https://www.home-assistant.io/docs/configuration/securing/)

---

## Next Steps

- [Configuration](Configuration) - Secure configuration options
- [Troubleshooting](Troubleshooting) - Debug security issues
- [Installation](Installation) - Secure installation guide
