package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type defaultRAMCollector struct{}

func (c *defaultRAMCollector) CollectRAM() (RAMMetrics, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return RAMMetrics{}, err
	}

	// Use Available instead of Free for more accurate representation
	// Available = Free + Buffers + Cached (reclaimable memory)
	// This makes Total ≈ Used + Available
	availableMB := v.Available / 1024 / 1024

	metrics := RAMMetrics{
		TotalMB: v.Total / 1024 / 1024,
		UsedMB:  v.Used / 1024 / 1024,
		FreeMB:  availableMB,
		Usage:   v.UsedPercent / 100.0,
	}

	return metrics, nil
}
