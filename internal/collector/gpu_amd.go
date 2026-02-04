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
	metrics := &GPUMetrics{
		Vendor: "amd",
		Model:  "AMD GPU", // Default
	}

	// Get GPU model name
	modelCmd := exec.Command("rocm-smi", "--showproductname")
	if modelOutput, err := modelCmd.Output(); err == nil {
		lines := strings.Split(string(modelOutput), "\n")
		for _, line := range lines {
			// Look for "GPU[0]          : Card Series:          AMD Radeon RX 6600 XT"
			if strings.Contains(line, "GPU[0]") && strings.Contains(line, "Card Series:") {
				// Extract card series name
				parts := strings.Split(line, "Card Series:")
				if len(parts) > 1 {
					metrics.Model = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Get usage, temperature, and VRAM info
	cmd := exec.Command("rocm-smi", "--showuse", "--showtemp", "--showmeminfo", "vram")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("rocm-smi failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// GPU use percentage: "GPU[0]          : GPU use (%): 0"
		if strings.Contains(line, "GPU use (%)") {
			re := regexp.MustCompile(`GPU use \(%\):\s*(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				util, _ := strconv.ParseFloat(matches[1], 64)
				metrics.Util = util / 100.0
			}
		}

		// Temperature (edge): "GPU[0]          : Temperature (Sensor edge) (C): 37.0"
		if strings.Contains(line, "Temperature (Sensor edge)") {
			re := regexp.MustCompile(`Temperature \(Sensor edge\) \(C\):\s*([\d.]+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				temp, _ := strconv.ParseFloat(matches[1], 64)
				metrics.TemperatureC = temp
			}
		}

		// VRAM Total (in Bytes): "GPU[0]          : VRAM Total Memory (B): 8573157376"
		if strings.Contains(line, "VRAM Total Memory (B)") {
			re := regexp.MustCompile(`VRAM Total Memory \(B\):\s*(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				vramBytes, _ := strconv.ParseUint(matches[1], 10, 64)
				metrics.VRAMTotalMB = vramBytes / 1024 / 1024 // Convert bytes to MB
			}
		}

		// VRAM Used (in Bytes): "GPU[0]          : VRAM Total Used Memory (B): 91881472"
		if strings.Contains(line, "VRAM Total Used Memory (B)") {
			re := regexp.MustCompile(`VRAM Total Used Memory \(B\):\s*(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				vramBytes, _ := strconv.ParseUint(matches[1], 10, 64)
				metrics.VRAMUsedMB = vramBytes / 1024 / 1024 // Convert bytes to MB
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
