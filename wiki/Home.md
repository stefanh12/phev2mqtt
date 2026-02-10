# phev2mqtt Wiki

Welcome to the **phev2mqtt** wiki! This is a comprehensive MQTT gateway for Mitsubishi Outlander PHEV vehicles with Home Assistant integration and advanced WiFi management features.

## Quick Links

### Getting Started
- **[Installation](Installation)** - Docker and Unraid setup instructions
- **[Configuration](Configuration)** - Complete configuration reference
- **[Quick Start Guide](Quick-Start)** - Get up and running in minutes

### Integration & Features
- **[Home Assistant Integration](Home-Assistant-Integration)** - Automatic discovery and setup
- **[MikroTik Integration](MikroTik-Integration)** - WiFi bridge management with RouterOS
- **[WiFi Management](WiFi-Management)** - Power saving and connection optimization

### Advanced Topics
- **[Development](Development)** - Contributing and local development setup
- **[Protocol Documentation](Protocol-Documentation)** - PHEV communication protocol details
- **[Security Best Practices](Security-Best-Practices)** - Securing your installation
- **[Troubleshooting](Troubleshooting)** - Common issues and solutions

## Project Overview

**phev2mqtt** is built on [buxtronix/phev2mqtt](https://github.com/buxtronix/phev2mqtt) and [CodeCutterUK/phev2mqtt](https://github.com/CodeCutterUK/phev2mqtt), optimized for running on Unraid with separate VLAN support and enhanced reliability through MikroTik WiFi client bridge integration.

### Key Features

✅ **Hot Configuration Reload** - Update settings without restarting the container  
✅ **Home Assistant Auto-Discovery** - Seamlessly integrates with Home Assistant MQTT  
✅ **Advanced WiFi Management** - Local and remote WiFi control with power saving modes  
✅ **MikroTik Integration** - Advanced monitoring and control via RouterOS  
✅ **Flexible Logging** - Multiple log levels for debugging and production use  
✅ **Security Hardened** - Input validation and secure credential handling  
✅ **Docker Ready** - Optimized for Docker and Unraid deployments

### Container Images

- **Stable**: `ghcr.io/stefanh12/phev2mqtt:latest`
- **Development**: `ghcr.io/stefanh12/phev2mqtt:master-dev`

## System Requirements

- Mitsubishi Outlander PHEV (MY18+)
- MQTT broker (e.g., Mosquitto)
- Docker or Unraid server
- WiFi bridge device (recommended: MikroTik RBSXTsq2nD)

## Compatibility

This gateway has been extensively tested with:
- **Vehicles**: MY20 and newer Outlander PHEV models
- **WiFi Bridges**: MikroTik RBSXTsq2nD (highly stable)
- **Home Assistant**: 2024.1 and newer
- **MQTT Brokers**: Mosquitto, HiveMQ

## Getting Help

- **Issues**: Report bugs on [GitHub Issues](../../issues)
- **Discussions**: Join the conversation in [GitHub Discussions](../../discussions)
- **Documentation**: Browse this wiki for comprehensive guides

## Contributing

Contributions are welcome! See the [Development](Development) guide for setup instructions and contribution guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
