# EMQX-PostgreSQL Bridge

A lightweight bridge service that connects EMQX MQTT broker to PostgreSQL database.

## Description

This service receives webhooks from EMQX and processes the data based on device type,
storing the results in a PostgreSQL database.

## Features

- Receives webhooks from EMQX
- Processes messages based on device type
- Extensible architecture for adding new device types
- Uses native SQL with pgx for optimal performance
- Lightweight and efficient
- Configuration via YAML file

## Requirements

- Go 1.24 or later
- PostgreSQL 13 or later
- EMQX with webhook capability

## Configuration

The service is configured via a YAML configuration file. By default, it looks for a config file in the following locations:

- `./config.yaml`
- `./config/config.yaml`
- `/etc/emqx-pg-bridge/config.yaml`
- `~/.emqx-pg-bridge/config.yaml`

You can also specify a custom config file location with the `-config` flag:

```bash
./emqx-pg-bridge -config /path/to/your/config.yaml
```



## Building

```bash
go build -o emqx-pg-bridge ./cmd/server
```

## Running

```bash
# Run with default config locations
./emqx-pg-bridge

# Run with specific config file
./emqx-pg-bridge -config ./my-config.yaml
```

## EMQX Configuration

Configure EMQX to send webhook requests to `http://your-service-host:8080/webhook`

## Topic Format

The topic format should follow: `homestay/{roomName}/{deviceName}/`

## User Properties

The service expects the following user properties in the MQTT message:

- `deviceType`: Type of device (e.g., "device-center", "normal")
- `roomId`: ID of the room
- `deviceId`: ID of the device
- `firmwareVersion`: Optional firmware version