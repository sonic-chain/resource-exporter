package device

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func GetGpu(gpu *NodeInfo) error {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("ERROR:: unable to initialize NVML: %d", ret)
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			return
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("ERROR:: unable to get device count: %d", ret)
	}
	gpu.Gpu.AttachedGpus = count

	driverVersion, ret := nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("ERROR:: unable to get device version: %d", ret)
	}
	gpu.Gpu.DriverVersion = driverVersion

	cudaDriverVersion, ret := nvml.SystemGetCudaDriverVersion_v2()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("ERROR:: unable to get device cuda version: %d", ret)
	}
	gpu.Gpu.CudaVersion = strconv.Itoa(cudaDriverVersion)

	for i := 0; i < count; i++ {
		var detail GpuDetail

		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			return fmt.Errorf("ERROR:: unable to get device index: %d", ret)
		}

		name, ret := device.GetName()
		if ret != nvml.SUCCESS {
			return fmt.Errorf("ERROR:: unable to get device name: %d", ret)
		}
		detail.ProductName = convertName(name)

		bar1MemoryInfo, ret := device.GetBAR1MemoryInfo()
		if ret != nvml.SUCCESS {
			return fmt.Errorf("ERROR:: unable to get bar1_memory of device at index %d: %d", i, ret)
		}

		detail.Bar1MemoryUsage.Total = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Total/1024/1024)
		detail.Bar1MemoryUsage.Used = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Used/1024/1024)
		detail.Bar1MemoryUsage.Free = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Free/1024/1024)

		memoryInfo, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			return fmt.Errorf("ERROR:: unable to get memory of device at index %d: %d", i, nvml.ErrorString(ret))
		}

		detail.FbMemoryUsage.Total = fmt.Sprintf("%d MiB", memoryInfo.Total/1024/1024)
		detail.FbMemoryUsage.Used = fmt.Sprintf("%d MiB", memoryInfo.Used/1024/1024)
		detail.FbMemoryUsage.Free = fmt.Sprintf("%d MiB", memoryInfo.Free/1024/1024)

		gpu.Gpu.Details = append(gpu.Gpu.Details, detail)
	}
	return nil
}

func convertName(name string) string {
	if strings.Contains(name, "Tesla") && !(strings.Contains(strings.ToUpper(name), "NVIDIA")) {
		return strings.Replace(name, "NVIDIA", "", 1)
	}

	if strings.Contains(name, "NVIDIA") {
		if strings.Contains(name, "Tesla") {
			return strings.Replace(name, "Tesla ", "", 1)
		}

		if strings.Contains(name, "GeForce") {
			name = strings.Replace(name, "GeForce ", "", 1)
		}

		if strings.Contains(name, "PCIE 40GB") {
			name = strings.Replace(name, "PCIE 40GB", "", 1)
		}
		if strings.Contains(name, "PCIE 80GB") {
			name = strings.Replace(name, "PCIE 80GB", "", 1)
		}

		return strings.Replace(name, "RTX ", "", 1)
	} else {
		if strings.Contains(name, "GeForce") {
			cpName := strings.Replace(name, "GeForce ", "NVIDIA", 1)
			return strings.Replace(cpName, "RTX", "", 1)
		}
	}
	return name
}
