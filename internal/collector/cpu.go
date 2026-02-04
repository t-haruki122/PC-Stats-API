package collector

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

type defaultCPUCollector struct{}

func (c *defaultCPUCollector) CollectCPU() (CPUMetrics, error) {
	metrics := CPUMetrics{
		Cores:   runtime.NumCPU(),
		Threads: runtime.NumCPU(),
	}

	// Get CPU model
	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		metrics.Model = info[0].ModelName
	}

	// Get CPU usage (average over 1 second)
	percentages, err := cpu.Percent(0, false)
	if err == nil && len(percentages) > 0 {
		metrics.Usage = percentages[0] / 100.0
	}

	// Get CPU frequency
	freq, err := cpu.Info()
	if err == nil && len(freq) > 0 {
		metrics.FrequencyMHz = freq[0].Mhz
	}

	// Get load average (Linux only)
	if runtime.GOOS == "linux" {
		loadAvg, err := load.Avg()
		if err == nil {
			metrics.LoadAvg = []float64{loadAvg.Load1, loadAvg.Load5, loadAvg.Load15}
		}
	}

	return metrics, nil
}
