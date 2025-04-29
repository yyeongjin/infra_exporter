
# infra_exporter

Deep within your systems, critical signals emerge.  
**infra_exporter** captures Linux security events and network activities — and delivers them precisely to Prometheus.

## ✨ Key Features
- **SSH Failed Login Detection**: Track intrusion attempts.
- **sudo Usage Monitoring**: Record privilege escalation events.
- **User Status Tracking**: Monitor UID/GID changes in real time.
- **Listening Ports Inspection**: Continuously track open ports.
- **External IP Connection Analysis**: Aggregate suspicious connections (Top N).
- **Sensitive File Monitoring**: Watch for presence of critical files like `/etc/passwd`.

## ⚡ Quick Start
```bash
go build -o infra_exporter ./cmd/infra_exporter
./infra_exporter --config ./config.yaml
```

You can also override settings easily via environment variables:
```bash
export MONITOR_SSH_FAILED=true
export MONITOR_SUDO_USAGE=true
./infra_exporter --config ./config.yaml
```

## 🌲 Directory Structure
```bash
infra_exporter/
├── cmd/
│   └── infra_exporter/           # Main application entry point
├── config/                       # Configuration files
│   └── config.go                 # Configuration logic
├── config.yaml                   # Configuration values
├── utils/                        # Utility functions
│   └── utils.go                  # System-related helper functions
└── README.md                     # This file
```

## 🔧 Configuration Example (config.yaml)
```yaml
monitor:
  ssh_failed: true
  sudo_usage: true
  user_status: true
  ports: true
  external_ip:
    enabled: true
    top_n: 5
  sensitive_file:
    enabled: true
    paths:
      - /etc/passwd
      - /etc/shadow
```
> Every field is overridable via environment variables.

## 📈 Metrics Overview
- **infra_ssh_failed_login_total**: Count of SSH failed login attempts
- **infra_sudo_usage_total**: Count of sudo usage events
- **infra_user_status**: User ID/GID information (Gauge)
- **infra_listening_ports_total**: Number of listening ports
- **infra_external_ip_connection_total**: External IP connection counts
- **infra_sensitive_file_status**: Sensitive file existence (1=exists, 0=missing)

## 🔍 PromQL Examples
- **SSH failed login rate over 5 minutes**:
```promql
rate(infra_ssh_failed_login_total[5m])
```
- **Sudo usage increase over 1 hour**:
```promql
increase(infra_sudo_usage_total[1h])
```
- **Current listening ports**:
```promql
infra_listening_ports_total
```
- **Top 5 external IP connections**:
```promql
topk(5, infra_external_ip_connection_total)
```
- **Check sensitive file existence**:
```promql
infra_sensitive_file_status{file="/etc/shadow"}
```

## 🚀 Endpoints
- `/metrics` - Exposes metrics for Prometheus scraping.

Default port is `:9101` (configurable).

## 🔥 Why infra_exporter?
- **Simple Configuration**: YAML and ENV support
- **Lightweight Execution**: Minimal resource footprint
- **Security-Focused Design**: Focus on the most critical events
- **Optimized for Prometheus**: Query-friendly metric naming

## 🛡️ License
Apache License 2.0
