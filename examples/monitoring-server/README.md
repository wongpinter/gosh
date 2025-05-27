# Monitoring Server SSH Example

This example demonstrates an SSH server that provides comprehensive system monitoring, metrics collection, and alerting capabilities.

## Features

- Real-time metrics collection
- Memory and runtime monitoring
- Request tracking and statistics
- Health checks and alerting
- Metrics export (JSON/CSV)
- Interactive dashboard
- Historical data tracking
- Background metric collection

## Setup

1. Generate SSH keys:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect and start monitoring:
   ```bash
   ssh -p 2228 monitor@localhost
   ```

## Available Commands

### Dashboard & Overview
- `dashboard` - Show comprehensive monitoring dashboard
- `health` - Perform system health check
- `alert` - Check for active alerts

### Metrics
- `metrics` - List all available metric types
- `metrics <type> [limit]` - Show specific metrics
- `memory` - Show detailed memory metrics
- `runtime` - Show Go runtime metrics
- `uptime` - Show uptime information
- `requests` - Show request statistics

### Monitoring Tools
- `watch <type>` - Watch specific metric type in real-time
- `export [format]` - Export metrics (json, csv)

### Utility
- `help` - Show all available commands

## Example Monitoring Session

```
$ ssh -p 2228 monitor@localhost
📊 Welcome to Monitoring Server!
Real-time system monitoring and metrics collection.
Type 'dashboard' for an overview or 'help' for commands.

monitor> dashboard
=== MONITORING DASHBOARD ===
Uptime: 5m30s
Memory: 2.1 MB / 8.7 MB
Goroutines: 12
Requests: 15
GC Runs: 3
Metrics Collected: 45
Last Updated: 15:30:45

monitor> memory
=== MEMORY METRICS ===
Allocated: 2.1 MB
System: 8.7 MB
Total Allocated: 5.2 MB
GC Runs: 3
Next GC: 4.2 MB
GC CPU Fraction: 0.0012

monitor> metrics memory 5
=== MEMORY METRICS (last 5) ===
[15:29:15] 2.05 MB
[15:29:45] 2.08 MB
[15:30:15] 2.10 MB
[15:30:45] 2.12 MB
[15:31:15] 2.15 MB

monitor> alert
✅ All systems normal - no alerts

monitor> health
=== HEALTH CHECK ===
Status: HEALTHY
Uptime: 6m15s
Memory Usage: 2.2 MB
Goroutines: 12
Last Check: 15:31:30

monitor> export json
=== METRICS EXPORT (JSON) ===
[
  {
    "timestamp": "2024-01-15T15:30:15Z",
    "type": "memory.alloc",
    "value": 2097152,
    "unit": "bytes"
  },
  {
    "timestamp": "2024-01-15T15:30:15Z",
    "type": "runtime.goroutines",
    "value": 12,
    "unit": "count"
  }
]
```

## Metrics Types

### Memory Metrics
- `memory.alloc` - Currently allocated memory
- `memory.sys` - System memory obtained from OS
- `memory.gc_runs` - Number of GC cycles

### Runtime Metrics
- `runtime.goroutines` - Number of active goroutines
- `runtime.cpus` - Number of CPUs
- `runtime.gc_runs` - Garbage collection runs

### Application Metrics
- `uptime.seconds` - Application uptime in seconds
- `requests.total` - Total number of requests
- `requests.rate` - Requests per second

## Alerting System

The monitoring server includes basic alerting for:

### Memory Alerts
- High memory usage (>100MB allocated)
- Memory leaks detection

### Performance Alerts
- High goroutine count (>100)
- High GC CPU usage (>10%)

### Custom Alerts
You can extend the alerting system by modifying the `checkAlerts()` method.

## Data Export

### JSON Export
```bash
monitor> export json
```
Exports all metrics in JSON format with timestamps, types, values, and tags.

### CSV Export
```bash
monitor> export csv
```
Exports metrics in CSV format suitable for spreadsheet analysis.

## Background Collection

The server automatically collects metrics every 30 seconds:
- Memory statistics
- Runtime information
- System performance data
- Application metrics

## Customization

### Adding Custom Metrics
```go
metricsCollector.AddMetric("custom.metric", value, "unit", tags)
```

### Modifying Collection Interval
```go
ticker := time.NewTicker(10 * time.Second) // Collect every 10 seconds
```

### Adding New Alert Rules
```go
if customCondition {
    alerts = append(alerts, "CUSTOM ALERT: description")
}
```

### Extending Metric Types
```go
case "custom":
    return h.getCustomMetrics(), 0
```

## Use Cases

### DevOps Monitoring
- Server health monitoring
- Performance tracking
- Resource usage analysis
- Capacity planning

### Application Monitoring
- Memory leak detection
- Performance bottlenecks
- Request rate monitoring
- Error tracking

### System Administration
- Real-time system status
- Historical trend analysis
- Alert management
- Automated reporting

### Development
- Performance profiling
- Resource optimization
- Load testing analysis
- Debug information

## Integration

The monitoring server can be integrated with:
- External monitoring systems
- Log aggregation platforms
- Alerting services
- Dashboard tools
- Automated deployment pipelines

## Security Considerations

- Read-only monitoring access
- No system modification capabilities
- Secure SSH authentication
- Audit logging of all commands
- Rate limiting for metric collection
