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

	metrics := RAMMetrics{
		TotalMB: v.Total / 1024 / 1024,
		UsedMB:  v.Used / 1024 / 1024,
		FreeMB:  v.Free / 1024 / 1024,
		Usage:   v.UsedPercent / 100.0,
	}

	return metrics, nil
}
