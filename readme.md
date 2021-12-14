# hkserver

Homekit server to bridge all non homekit device

## Configuration

Configuration example

```yaml
logLevel: debug

homekit:
  enable: true
  name: 'homeb-dev'
  pin: '00102003'
  storagePath: './db'
  port: "48000"

http:
  port: 8080

mqtt:
  brokers:
    - "tcp://192.168.1.40:1883"
  username: mqtt
  password: mqtt

discovery:
  enable: true

devices:
  - name: "my device"
    ip: 192.168.20.54
    type: shelly1
    id: 1 # optional but recommended: unique id of device
```

### Devices

**shelly 1 v3**

```yaml
name: "my device"
ip: 192.168.20.54
type: shelly1
```

**shelly 1pm**

```yaml
name: "my device"
ip: 192.168.20.54
type: shelly1pm
```

**daitem allarm**

```yaml
name: "my device"
ip: 192.168.20.54
type: daitem
options:
  email: "..."
  password: "..."
  masterCode: "..."
```

## Startup

```bash
# Default
hkserver 

# With custom configuratioin path
hkserver --config <config_path.yml>
```

### Service

To use hkserver as a service

```bash
cd /etc/systemd/system
vi hkserver.service

# copy scripts/hkserver.service

systemctl start hkserver
systemctl enable hkserver
```
