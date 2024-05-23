package device

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func GetGpu(gpu *NodeInfo) nvml.Return {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return ret
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			return
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return ret
	}
	gpu.Gpu.AttachedGpus = count

	driverVersion, ret := nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		return ret
	}
	gpu.Gpu.DriverVersion = driverVersion

	cudaDriverVersion, ret := nvml.SystemGetCudaDriverVersion_v2()
	if ret != nvml.SUCCESS {
		return ret
	}
	gpu.Gpu.CudaVersion = strconv.Itoa(cudaDriverVersion)

	for i := 0; i < count; i++ {
		var detail GpuDetail

		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret == nvml.SUCCESS {
			name, ret := device.GetName()
			if ret == nvml.SUCCESS {
				detail.ProductName = strings.ToUpper(convertName(name))
			}
			bar1MemoryInfo, ret := device.GetBAR1MemoryInfo()
			if ret == nvml.SUCCESS {
				detail.Bar1MemoryUsage.Total = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Total/1024/1024)
				detail.Bar1MemoryUsage.Used = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Used/1024/1024)
				detail.Bar1MemoryUsage.Free = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Free/1024/1024)
			}
			memoryInfo, ret := device.GetMemoryInfo()
			if ret == nvml.SUCCESS {
				detail.FbMemoryUsage.Total = fmt.Sprintf("%d MiB", memoryInfo.Total/1024/1024)
				detail.FbMemoryUsage.Used = fmt.Sprintf("%d MiB", memoryInfo.Used/1024/1024)
				detail.FbMemoryUsage.Free = fmt.Sprintf("%d MiB", memoryInfo.Free/1024/1024)
			}

			processes, err := deviceGetAllRunningProcesses(device)
			if err == nil {
				if len(processes) > 0 {
					detail.Status = "occupied"
				} else {
					detail.Status = "available"
				}
			}
			gpu.Gpu.Details = append(gpu.Gpu.Details, detail)
		}
	}
	return ret
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
		return strings.Replace(name, "RTX ", "", 1)
	} else {
		if strings.Contains(name, "GeForce") {
			cpName := strings.Replace(name, "GeForce ", "NVIDIA", 1)
			return strings.Replace(cpName, "RTX", "", 1)
		}
	}
	return name
}

func deviceGetAllRunningProcesses(device nvml.Device) ([]nvml.ProcessInfo, error) {
	process1, ret := device.GetComputeRunningProcesses()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("ERROR:: unable to get device index: %d", ret)
	}

	process2, ret := device.GetGraphicsRunningProcesses()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("ERROR:: unable to get device index: %d", ret)
	}

	processInfo := make([]nvml.ProcessInfo, 0)

	if len(process1) > 0 {
		processInfo = append(processInfo, process1...)
	}
	if len(process2) > 0 {
		processInfo = append(processInfo, process2...)
	}
	return processInfo, nil
}
