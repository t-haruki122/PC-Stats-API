package collector

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type nvidiaCollector struct{}

func (n *nvidiaCollector) Vendor() string {
	return "nvidia"
}

func (n *nvidiaCollector) CollectGPU() (*GPUMetrics, error) {
	// nvidia-smi --query-gpu=name,utilization.gpu,temperature.gpu,memory.total,memory.used --format=csv,noheader,nounits
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=name,utilization.gpu,temperature.gpu,memory.total,memory.used",
		"--format=csv,noheader,nounits")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi failed: %w", err)
	}

	// Parse CSV output
	reader := csv.NewReader(strings.NewReader(string(output)))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil || len(records) == 0 {
		return nil, fmt.Errorf("failed to parse nvidia-smi output")
	}

	// Use first GPU
	record := records[0]
	if len(record) < 5 {
		return nil, fmt.Errorf("unexpected nvidia-smi output format")
	}

	util, _ := strconv.ParseFloat(record[1], 64)
	temp, _ := strconv.ParseFloat(record[2], 64)
	vramTotal, _ := strconv.ParseUint(record[3], 10, 64)
	vramUsed, _ := strconv.ParseUint(record[4], 10, 64)

	return &GPUMetrics{
		Vendor:       "nvidia",
		Model:        strings.TrimSpace(record[0]),
		Util:         util / 100.0,
		TemperatureC: temp,
		VRAMTotalMB:  vramTotal,
		VRAMUsedMB:   vramUsed,
	}, nil
}
