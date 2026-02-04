package collector

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type amdCollector struct{}

func (a *amdCollector) Vendor() string {
	return "amd"
}

func (a *amdCollector) CollectGPU() (*GPUMetrics, error) {
	if runtime.GOOS == "linux" {
		return a.collectLinux()
	}
	return a.collectWindows()
}

func (a *amdCollector) collectLinux() (*GPUMetrics, error) {
	// rocm-smi --showuse --showtemp --showmeminfo vram
	cmd := exec.Command("rocm-smi", "--showuse", "--showtemp", "--showmeminfo", "vram")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("rocm-smi failed: %w", err)
	}

	metrics := &GPUMetrics{
		Vendor: "amd",
	}

	lines := strings.Split(string(output), "\n")

	// Parse output (simplified - actual parsing depends on rocm-smi format)
	for _, line := range lines {
		// GPU use percentage
		if strings.Contains(line, "GPU use") {
			re := regexp.MustCompile(`(\d+)%`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				util, _ := strconv.ParseFloat(matches[1], 64)
				metrics.Util = util / 100.0
			}
		}
		// Temperature
		if strings.Contains(line, "Temperature") {
			re := regexp.MustCompile(`(\d+\.?\d*)c`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				temp, _ := strconv.ParseFloat(matches[1], 64)
				metrics.TemperatureC = temp
			}
		}
		// VRAM
		if strings.Contains(line, "VRAM Total") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				vram, _ := strconv.ParseUint(matches[1], 10, 64)
				metrics.VRAMTotalMB = vram
			}
		}
		if strings.Contains(line, "VRAM Used") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				vram, _ := strconv.ParseUint(matches[1], 10, 64)
				metrics.VRAMUsedMB = vram
			}
		}
	}

	return metrics, nil
}

func (a *amdCollector) collectWindows() (*GPUMetrics, error) {
	// Simplified Windows implementation using WMI
	// Get GPU name and basic info
	cmd := exec.Command("powershell", "-Command",
		"Get-WmiObject Win32_VideoController | Select-Object Name, AdapterRAM | ConvertTo-Json")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("WMI query failed: %w", err)
	}

	// Basic parsing (simplified)
	metrics := &GPUMetrics{
		Vendor: "amd",
		Model:  "AMD GPU", // Default
		Util:   0.0,       // Not easily available via WMI
	}

	// Try to extract model name
	if strings.Contains(string(output), "Name") {
		re := regexp.MustCompile(`"Name":\s*"([^"]+)"`)
		if matches := re.FindStringSubmatch(string(output)); len(matches) > 1 {
			metrics.Model = matches[1]
		}
	}

	return metrics, nil
}
