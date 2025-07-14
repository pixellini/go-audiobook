package coqui

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Device represents the compute device for TTS synthesis.
// Supported devices include CPU, CUDA (NVIDIA GPU), MPS (Apple Silicon), and auto-detection.
type Device string

// Runtime represents system runtime characteristics.
// Used internally for device detection and platform-specific optimizations.
type Runtime string

const (
	// Auto enables automatic device detection based on available hardware.
	Auto Device = "auto"
	// CPU forces CPU-only synthesis (slowest but most compatible).
	CPU Device = "cpu"
	// CUDA enables NVIDIA GPU acceleration (requires CUDA installation).
	CUDA Device = "cuda"
	// MPS enables Apple Silicon GPU acceleration (macOS only).
	MPS Device = "mps"
)

// String returns the string representation of the Device.
func (d Device) String() string {
	return string(d)
}

// IsValid checks if the device type is supported.
// Returns true for Auto, CPU, CUDA, and MPS devices.
func (d Device) IsValid() bool {
	switch d {
	case Auto, CPU, CUDA, MPS:
		return true
	default:
		return false
	}
}

const (
	// Darwin represents the macOS operating system.
	Darwin Runtime = "darwin"
	// ARM64 represents ARM64 processor architecture (Apple Silicon).
	ARM64 Runtime = "arm64"
)

// String returns the string representation of the Runtime.
func (r Runtime) String() string {
	return string(r)
}

// DetectDevice automatically selects the best available compute device.
// Priority order: CUDA (if available) > MPS (macOS ARM64) > CPU (fallback).
func DetectDevice(d Device) Device {
	if isCudaAvailable() {
		fmt.Println("CUDA available, using CUDA device.")
		return CUDA
	}

	if runtime.GOOS == string(Darwin) && runtime.GOARCH == string(ARM64) {
		fmt.Println("ARM64 architecture detected, using MPS.")
		return MPS
	}

	// Default to "cpu" if no other device is available
	fmt.Println("No GPU detected or CUDA not available, using CPU device.")
	return CPU
}

// isCudaAvailable checks if NVIDIA GPU and drivers are available.
// This performs a basic check by looking for the nvidia-smi command.
// Note: This does not guarantee CUDA is properly configured for TTS.
func isCudaAvailable() bool {
	_, err := exec.LookPath("nvidia-smi")
	return err == nil
}
