package device

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func GetGpu(gpu *NodeInfo) {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Printf("Unable to get device count: %v", nvml.ErrorString(ret))
		return
	}
	gpu.Gpu.AttachedGpus = count

	driverVersion, ret := nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		log.Printf("Unable to get device count: %v", nvml.ErrorString(ret))
		return
	}
	gpu.Gpu.DriverVersion = driverVersion

	cudaDriverVersion, ret := nvml.SystemGetCudaDriverVersion_v2()
	if ret != nvml.SUCCESS {
		log.Printf("Unable to get device count: %v", nvml.ErrorString(ret))
		return
	}
	gpu.Gpu.CudaVersion = strconv.Itoa(cudaDriverVersion)

	for i := 0; i < count; i++ {
		var detail GpuDetail

		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Printf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
			return
		}

		name, ret := device.GetName()
		if ret != nvml.SUCCESS {
			log.Printf("Unable to get name of device at index %d: %v", i, nvml.ErrorString(ret))
			return
		}
		detail.ProductName = convertName(name)

		bar1MemoryInfo, ret := device.GetBAR1MemoryInfo()
		if ret != nvml.SUCCESS {
			log.Printf("Unable to get bar1_memory of device at index %d: %v", i, nvml.ErrorString(ret))
			return
		}

		detail.Bar1MemoryUsage.Total = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Total/1024/1024)
		detail.Bar1MemoryUsage.Used = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Used/1024/1024)
		detail.Bar1MemoryUsage.Free = fmt.Sprintf("%d MiB", bar1MemoryInfo.Bar1Free/1024/1024)

		memoryInfo, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			log.Printf("Unable to get memory of device at index %d: %v", i, nvml.ErrorString(ret))
			return
		}

		detail.FbMemoryUsage.Total = fmt.Sprintf("%d MiB", memoryInfo.Total/1024/1024)
		detail.FbMemoryUsage.Used = fmt.Sprintf("%d MiB", memoryInfo.Used/1024/1024)
		detail.FbMemoryUsage.Free = fmt.Sprintf("%d MiB", memoryInfo.Free/1024/1024)

		gpu.Gpu.Details = append(gpu.Gpu.Details, detail)
	}
}

func convertName(name string) string {
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
