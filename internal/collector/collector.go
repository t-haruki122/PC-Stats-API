package collector

import "time"

// Collector is the interface for collecting metrics
type Collector interface {
	Collect() (*MetricSample, error)
}

// MetricSample represents a complete snapshot of system metrics
type MetricSample struct {
	Timestamp time.Time   `json:"timestamp"`
	CPU       CPUMetrics  `json:"cpu"`
	RAM       RAMMetrics  `json:"ram"`
	GPU       *GPUMetrics `json:"gpu,omitempty"`
}

// CPUMetrics contains CPU information
type CPUMetrics struct {
	Model        string    `json:"model"`
	Cores        int       `json:"cores"`
	Threads      int       `json:"threads"`
	Usage        float64   `json:"usage"`              // 0.0-1.0
	LoadAvg      []float64 `json:"load_avg,omitempty"` // Linux only
	FrequencyMHz float64   `json:"frequency_mhz,omitempty"`
}

// RAMMetrics contains RAM information
type RAMMetrics struct {
	TotalMB uint64  `json:"total_mb"`
	UsedMB  uint64  `json:"used_mb"`
	FreeMB  uint64  `json:"free_mb"`
	Usage   float64 `json:"usage"` // 0.0-1.0
}

// GPUMetrics contains GPU information
type GPUMetrics struct {
	Vendor       string  `json:"vendor"` // "nvidia" or "amd"
	Model        string  `json:"model"`
	Util         float64 `json:"util"` // 0.0-1.0
	TemperatureC float64 `json:"temperature_c,omitempty"`
	VRAMTotalMB  uint64  `json:"vram_total_mb,omitempty"`
	VRAMUsedMB   uint64  `json:"vram_used_mb,omitempty"`
}

// SystemCollector collects all system metrics
type SystemCollector struct {
	cpuCollector CPUCollector
	ramCollector RAMCollector
	gpuCollector GPUCollector
}

// NewSystemCollector creates a new system collector
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		cpuCollector: &defaultCPUCollector{},
		ramCollector: &defaultRAMCollector{},
		gpuCollector: detectGPU(),
	}
}

// Collect gathers all metrics
func (sc *SystemCollector) Collect() (*MetricSample, error) {
	sample := &MetricSample{
		Timestamp: time.Now(),
	}

	// Collect CPU
	cpu, err := sc.cpuCollector.CollectCPU()
	if err != nil {
		return nil, err
	}
	sample.CPU = cpu

	// Collect RAM
	ram, err := sc.ramCollector.CollectRAM()
	if err != nil {
		return nil, err
	}
	sample.RAM = ram

	// Collect GPU (optional)
	if sc.gpuCollector != nil {
		gpu, err := sc.gpuCollector.CollectGPU()
		if err == nil && gpu != nil {
			sample.GPU = gpu
		}
	}

	return sample, nil
}

// CPUCollector interface
type CPUCollector interface {
	CollectCPU() (CPUMetrics, error)
}

// RAMCollector interface
type RAMCollector interface {
	CollectRAM() (RAMMetrics, error)
}

// GPUCollector interface
type GPUCollector interface {
	CollectGPU() (*GPUMetrics, error)
	Vendor() string
}
