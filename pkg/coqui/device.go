package coqui

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Devices constants
type Device string
type Runtime string

const (
	Auto Device = "auto"
	CPU  Device = "cpu"
	CUDA Device = "cuda"
	MPS  Device = "mps"
)

func (d Device) String() string {
	return string(d)
}

// IsValid checks if the model is supported
func (d Device) IsValid() bool {
	switch d {
	case Auto, CPU, CUDA, MPS:
		return true
	default:
		return false
	}
}

const (
	// Architecture and OS constants
	Darwin    Runtime = "darwin"
	ArchARM64 Runtime = "arm64"
)

func (r Runtime) String() string {
	return string(r)
}

func DetectDevice(d Device) Device {
	if isCudaAvailable() {
		fmt.Println("CUDA available, using CUDA device.")
		return CUDA
	}

	if runtime.GOOS == string(Darwin) && runtime.GOARCH == string(ArchARM64) {
		fmt.Println("ARM64 architecture detected, using MPS.")
		return MPS
	}

	// Default to "cpu" if no other device is available
	fmt.Println("No GPU detected or CUDA not available, using CPU device.")
	return CPU
}

// This only checks if the system has the "nvidia-smi" command available,
// but does not guarantee that CUDA is properly installed or configured.
// If CUDA is not properly configured, the user will likely see an error when starting the TTS process.
// TODO: Check for Go packages that can verify this better.
func isCudaAvailable() bool {
	_, err := exec.LookPath("nvidia-smi")
	return err == nil
}
