package device

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

// Devices constants
type Device string

const (
	DeviceAuto Device = "auto"
	DeviceCPU  Device = "cpu"
	DeviceCUDA Device = "cuda"
	DeviceMPS  Device = "mps"
)

// Architecture and OS constants
type DeviceOS string

const (
	osDarwin  DeviceOS = "darwin"
	archARM64 DeviceOS = "arm64"
)

// Singleton instance for the user's device.
type DeviceManager struct {
	Device Device
	set    bool
}

var Manager = &DeviceManager{
	Device: DeviceAuto,
	set:    false,
}

func (dm *DeviceManager) Init() Device {
	if dm.set {
		return dm.Device
	}

	dm.Device = assignDevice()
	dm.set = true
	return dm.Device
}

// This only checks if the system has the "nvidia-smi" command available,
// but does not guarantee that CUDA is properly installed or configured.
// If CUDA is not properly configured, the user will likely see an error when starting the TTS process.
// TODO: Check for Go packages that can verify this better.
func cudaIsAvailable() bool {
	_, err := exec.LookPath("nvidia-smi")
	return err == nil
}

func assignDevice() Device {
	device := Device(viper.GetString("tts.device"))

	// User has specified a device
	if device != "" && device != DeviceAuto {
		return device
	}

	if cudaIsAvailable() {
		fmt.Println("CUDA available, using CUDA device.")
		return DeviceCUDA
	}

	if runtime.GOOS == string(osDarwin) && runtime.GOARCH == string(archARM64) {
		fmt.Println("ARM64 architecture detected, using MPS.")
		return DeviceMPS
	}

	// Default to "cpu" if no other device is available
	fmt.Println("No GPU detected or CUDA not available, using CPU device.")
	return DeviceCPU
}
