package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"repo.nusatek.id/sugeng/gosh"
)

// MetricsCollector collects and stores system metrics
type MetricsCollector struct {
	startTime time.Time
	metrics   []Metric
	maxMetrics int
}

// Metric represents a single metric data point
type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime:  time.Now(),
		metrics:    make([]Metric, 0),
		maxMetrics: 1000,
	}
}

// AddMetric adds a metric to the collection
func (mc *MetricsCollector) AddMetric(metricType string, value float64, unit string, tags map[string]string) {
	metric := Metric{
		Timestamp: time.Now(),
		Type:      metricType,
		Value:     value,
		Unit:      unit,
		Tags:      tags,
	}
	
	mc.metrics = append(mc.metrics, metric)
	
	// Keep only the last maxMetrics
	if len(mc.metrics) > mc.maxMetrics {
		mc.metrics = mc.metrics[len(mc.metrics)-mc.maxMetrics:]
	}
}

// GetMetrics returns metrics of a specific type
func (mc *MetricsCollector) GetMetrics(metricType string, limit int) []Metric {
	var filtered []Metric
	
	for i := len(mc.metrics) - 1; i >= 0 && len(filtered) < limit; i-- {
		if mc.metrics[i].Type == metricType {
			filtered = append(filtered, mc.metrics[i])
		}
	}
	
	// Reverse to get chronological order
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}
	
	return filtered
}

// Global metrics collector
var metricsCollector = NewMetricsCollector()

// MonitoringHandler implements monitoring and metrics commands
type MonitoringHandler struct {
	startTime time.Time
	requests  int
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler() *MonitoringHandler {
	handler := &MonitoringHandler{
		startTime: time.Now(),
		requests:  0,
	}
	
	// Start background metric collection
	go handler.collectMetrics()
	
	return handler
}

// Execute implements the CommandHandler interface
func (h *MonitoringHandler) Execute(cmd string) (string, uint32) {
	h.requests++
	
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return "", 0
	}
	
	command := parts[0]
	args := parts[1:]
	
	switch command {
	case "metrics":
		return h.getMetrics(args)
	case "memory":
		return h.getMemoryMetrics(), 0
	case "runtime":
		return h.getRuntimeMetrics(), 0
	case "uptime":
		return h.getUptimeMetrics(), 0
	case "requests":
		return h.getRequestMetrics(), 0
	case "dashboard":
		return h.getDashboard(), 0
	case "export":
		return h.exportMetrics(args)
	case "alert":
		return h.checkAlerts(), 0
	case "health":
		return h.getHealthCheck(), 0
	case "watch":
		return h.watchMetrics(args)
	case "help":
		return h.getHelp(), 0
	default:
		return fmt.Sprintf("Unknown command: %s\nType 'help' for available commands", command), 1
	}
}

func (h *MonitoringHandler) getMetrics(args []string) (string, uint32) {
	if len(args) == 0 {
		return h.listMetricTypes(), 0
	}
	
	metricType := args[0]
	limit := 10
	
	if len(args) > 1 {
		if l, err := strconv.Atoi(args[1]); err == nil && l > 0 {
			limit = l
		}
	}
	
	metrics := metricsCollector.GetMetrics(metricType, limit)
	if len(metrics) == 0 {
		return fmt.Sprintf("No metrics found for type: %s", metricType), 0
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("=== %s METRICS (last %d) ===\n", strings.ToUpper(metricType), len(metrics)))
	
	for _, metric := range metrics {
		timestamp := metric.Timestamp.Format("15:04:05")
		result.WriteString(fmt.Sprintf("[%s] %.2f %s", timestamp, metric.Value, metric.Unit))
		
		if len(metric.Tags) > 0 {
			result.WriteString(" (")
			first := true
			for k, v := range metric.Tags {
				if !first {
					result.WriteString(", ")
				}
				result.WriteString(fmt.Sprintf("%s=%s", k, v))
				first = false
			}
			result.WriteString(")")
		}
		result.WriteString("\n")
	}
	
	return result.String(), 0
}

func (h *MonitoringHandler) listMetricTypes() string {
	types := make(map[string]int)
	
	for _, metric := range metricsCollector.metrics {
		types[metric.Type]++
	}
	
	var result strings.Builder
	result.WriteString("=== AVAILABLE METRICS ===\n")
	
	if len(types) == 0 {
		result.WriteString("No metrics collected yet.\n")
		return result.String()
	}
	
	for metricType, count := range types {
		result.WriteString(fmt.Sprintf("- %s (%d data points)\n", metricType, count))
	}
	
	result.WriteString("\nUsage: metrics <type> [limit]\n")
	result.WriteString("Example: metrics memory 20\n")
	
	return result.String()
}

func (h *MonitoringHandler) getMemoryMetrics() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Add current metrics
	metricsCollector.AddMetric("memory.alloc", float64(m.Alloc), "bytes", nil)
	metricsCollector.AddMetric("memory.sys", float64(m.Sys), "bytes", nil)
	metricsCollector.AddMetric("memory.gc_runs", float64(m.NumGC), "count", nil)
	
	return fmt.Sprintf("=== MEMORY METRICS ===\n"+
		"Allocated: %s\n"+
		"System: %s\n"+
		"Total Allocated: %s\n"+
		"GC Runs: %d\n"+
		"Next GC: %s\n"+
		"GC CPU Fraction: %.4f",
		h.formatBytes(m.Alloc),
		h.formatBytes(m.Sys),
		h.formatBytes(m.TotalAlloc),
		m.NumGC,
		h.formatBytes(m.NextGC),
		m.GCCPUFraction)
}

func (h *MonitoringHandler) getRuntimeMetrics() string {
	// Add current metrics
	metricsCollector.AddMetric("runtime.goroutines", float64(runtime.NumGoroutine()), "count", nil)
	metricsCollector.AddMetric("runtime.cpus", float64(runtime.NumCPU()), "count", nil)
	
	return fmt.Sprintf("=== RUNTIME METRICS ===\n"+
		"Go Version: %s\n"+
		"OS/Arch: %s/%s\n"+
		"CPUs: %d\n"+
		"Goroutines: %d\n"+
		"CGO Calls: %d",
		runtime.Version(),
		runtime.GOOS, runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		runtime.NumCgoCall())
}

func (h *MonitoringHandler) getUptimeMetrics() string {
	uptime := time.Since(h.startTime)
	serverUptime := time.Since(metricsCollector.startTime)
	
	// Add uptime metric
	metricsCollector.AddMetric("uptime.seconds", uptime.Seconds(), "seconds", map[string]string{"component": "handler"})
	
	return fmt.Sprintf("=== UPTIME METRICS ===\n"+
		"Handler Uptime: %v\n"+
		"Server Uptime: %v\n"+
		"Started: %s",
		uptime.Round(time.Second),
		serverUptime.Round(time.Second),
		h.startTime.Format("2006-01-02 15:04:05"))
}

func (h *MonitoringHandler) getRequestMetrics() string {
	uptime := time.Since(h.startTime)
	requestsPerSecond := float64(h.requests) / uptime.Seconds()
	
	// Add request metrics
	metricsCollector.AddMetric("requests.total", float64(h.requests), "count", nil)
	metricsCollector.AddMetric("requests.rate", requestsPerSecond, "req/sec", nil)
	
	return fmt.Sprintf("=== REQUEST METRICS ===\n"+
		"Total Requests: %d\n"+
		"Requests/Second: %.2f\n"+
		"Average Response Time: N/A",
		h.requests, requestsPerSecond)
}

func (h *MonitoringHandler) getDashboard() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	uptime := time.Since(h.startTime)
	
	return fmt.Sprintf("=== MONITORING DASHBOARD ===\n"+
		"Uptime: %v\n"+
		"Memory: %s / %s\n"+
		"Goroutines: %d\n"+
		"Requests: %d\n"+
		"GC Runs: %d\n"+
		"Metrics Collected: %d\n"+
		"Last Updated: %s",
		uptime.Round(time.Second),
		h.formatBytes(m.Alloc), h.formatBytes(m.Sys),
		runtime.NumGoroutine(),
		h.requests,
		m.NumGC,
		len(metricsCollector.metrics),
		time.Now().Format("15:04:05"))
}

func (h *MonitoringHandler) exportMetrics(args []string) (string, uint32) {
	format := "json"
	if len(args) > 0 {
		format = args[0]
	}
	
	switch format {
	case "json":
		data, err := json.MarshalIndent(metricsCollector.metrics, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error exporting metrics: %v", err), 1
		}
		return fmt.Sprintf("=== METRICS EXPORT (JSON) ===\n%s", string(data)), 0
	case "csv":
		return h.exportCSV(), 0
	default:
		return "Supported formats: json, csv", 1
	}
}

func (h *MonitoringHandler) exportCSV() string {
	var result strings.Builder
	result.WriteString("=== METRICS EXPORT (CSV) ===\n")
	result.WriteString("timestamp,type,value,unit,tags\n")
	
	for _, metric := range metricsCollector.metrics {
		tags := ""
		if len(metric.Tags) > 0 {
			tagPairs := make([]string, 0, len(metric.Tags))
			for k, v := range metric.Tags {
				tagPairs = append(tagPairs, fmt.Sprintf("%s=%s", k, v))
			}
			tags = strings.Join(tagPairs, ";")
		}
		
		result.WriteString(fmt.Sprintf("%s,%s,%.2f,%s,%s\n",
			metric.Timestamp.Format("2006-01-02T15:04:05Z"),
			metric.Type,
			metric.Value,
			metric.Unit,
			tags))
	}
	
	return result.String()
}

func (h *MonitoringHandler) checkAlerts() string {
	var alerts []string
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Memory alerts
	if m.Alloc > 100*1024*1024 { // 100MB
		alerts = append(alerts, fmt.Sprintf("HIGH MEMORY: %s allocated", h.formatBytes(m.Alloc)))
	}
	
	// Goroutine alerts
	if runtime.NumGoroutine() > 100 {
		alerts = append(alerts, fmt.Sprintf("HIGH GOROUTINES: %d active", runtime.NumGoroutine()))
	}
	
	// GC alerts
	if m.GCCPUFraction > 0.1 {
		alerts = append(alerts, fmt.Sprintf("HIGH GC CPU: %.2f%% CPU time", m.GCCPUFraction*100))
	}
	
	if len(alerts) == 0 {
		return "✅ All systems normal - no alerts"
	}
	
	result := "🚨 ACTIVE ALERTS:\n"
	for i, alert := range alerts {
		result += fmt.Sprintf("%d. %s\n", i+1, alert)
	}
	
	return result
}

func (h *MonitoringHandler) getHealthCheck() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	uptime := time.Since(h.startTime)
	
	status := "HEALTHY"
	if m.Alloc > 100*1024*1024 || runtime.NumGoroutine() > 100 {
		status = "WARNING"
	}
	
	return fmt.Sprintf("=== HEALTH CHECK ===\n"+
		"Status: %s\n"+
		"Uptime: %v\n"+
		"Memory Usage: %s\n"+
		"Goroutines: %d\n"+
		"Last Check: %s",
		status,
		uptime.Round(time.Second),
		h.formatBytes(m.Alloc),
		runtime.NumGoroutine(),
		time.Now().Format("15:04:05"))
}

func (h *MonitoringHandler) watchMetrics(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: watch <metric_type>\nExample: watch memory", 1
	}
	
	metricType := args[0]
	recent := metricsCollector.GetMetrics(metricType, 5)
	
	if len(recent) == 0 {
		return fmt.Sprintf("No recent data for metric: %s", metricType), 1
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("=== WATCHING %s ===\n", strings.ToUpper(metricType)))
	
	for _, metric := range recent {
		result.WriteString(fmt.Sprintf("[%s] %.2f %s\n",
			metric.Timestamp.Format("15:04:05"),
			metric.Value,
			metric.Unit))
	}
	
	result.WriteString("\nRefresh with: watch " + metricType)
	return result.String(), 0
}

func (h *MonitoringHandler) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// Collect various metrics
		metricsCollector.AddMetric("memory.alloc", float64(m.Alloc), "bytes", nil)
		metricsCollector.AddMetric("memory.sys", float64(m.Sys), "bytes", nil)
		metricsCollector.AddMetric("runtime.goroutines", float64(runtime.NumGoroutine()), "count", nil)
		metricsCollector.AddMetric("runtime.gc_runs", float64(m.NumGC), "count", nil)
		metricsCollector.AddMetric("uptime.seconds", time.Since(h.startTime).Seconds(), "seconds", nil)
	}
}

func (h *MonitoringHandler) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (h *MonitoringHandler) getHelp() string {
	return `Monitoring Commands:
- dashboard          Show monitoring dashboard
- metrics [type]     List metrics or show specific type
- memory             Show memory metrics
- runtime            Show runtime metrics
- uptime             Show uptime information
- requests           Show request statistics
- health             Perform health check
- alert              Check for alerts
- watch <type>       Watch specific metric type
- export [format]    Export metrics (json, csv)
- help               Show this help

Examples:
- metrics memory 20  Show last 20 memory metrics
- watch runtime      Watch runtime metrics
- export json        Export all metrics as JSON`
}

// GetPrompt implements the CommandHandler interface
func (h *MonitoringHandler) GetPrompt() string {
	return "monitor> "
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *MonitoringHandler) GetWelcomeMessage() string {
	return "📊 Welcome to Monitoring Server!\n" +
		"Real-time system monitoring and metrics collection.\n" +
		"Type 'dashboard' for an overview or 'help' for commands."
}

func main() {
	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2228"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "monitoring_server.log"

	// Create monitoring handler
	handler := NewMonitoringHandler()

	// Create and start server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Monitoring Server started on port 2228!")
	log.Println("Connect with: ssh -p 2228 monitor@localhost")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down monitoring server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
