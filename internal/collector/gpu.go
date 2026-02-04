package collector

import (
	"os/exec"
)

// detectGPU attempts to detect available GPU
func detectGPU() GPUCollector {
	// Try NVIDIA first
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		return &nvidiaCollector{}
	}

	// Try AMD on Linux
	if _, err := exec.LookPath("rocm-smi"); err == nil {
		return &amdCollector{}
	}

	// Try AMD on Windows (WMI-based)
	// For now, we'll implement a basic version
	// TODO: Implement AMD Windows support

	return nil // No GPU detected
}
