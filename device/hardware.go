package device

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"math"
	"strings"
	"time"
)

func GetHardwareData(node *NodeInfo) error {

	diskStat, err := disk.Usage("/")
	if err != nil {
		return err
	}
	node.Storage.Total = fmt.Sprintf("%d GiB", diskStat.Total/1024/1024/1024)
	node.Storage.Used = fmt.Sprintf("%d GiB", diskStat.Used/1024/1024/1024)
	node.Storage.Free = fmt.Sprintf("%d GiB", diskStat.Free/1024/1024/1024)

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	node.Memory.Total = fmt.Sprintf("%d GiB", vmStat.Total/1024/1024/1024)
	node.Memory.Used = fmt.Sprintf("%d GiB", vmStat.Used/1024/1024/1024)
	node.Memory.Free = fmt.Sprintf("%d GiB", vmStat.Available/1024/1024/1024)

	counts, err := cpu.Counts(true)
	if err != nil {
		return err
	}

	infoStats, err := cpu.Info()
	if err != nil {
		return err
	}

	percent, err := cpu.Percent(5*time.Second, false)
	if err != nil {
		return err
	}

	var useCount int
	if percent[0] > 1.0 {
		useCount = int(math.Round(float64(counts) / 1.6))
	} else {
		useCount = int(math.Round(float64(counts) * percent[0]))
	}

	node.Cpu.Total = fmt.Sprintf("%d", counts)
	node.Cpu.Used = fmt.Sprintf("%d", useCount)
	node.Cpu.Free = fmt.Sprintf("%d", counts-useCount)

	node.Vcpu.Total = node.Cpu.Total
	node.Vcpu.Used = node.Cpu.Used
	node.Vcpu.Free = node.Cpu.Free

	hostStat, err := host.Info()
	if err != nil {
		return err
	}

	node.MachineId = hostStat.HostID
	if strings.Contains(strings.ToLower(infoStats[0].ModelName), "intel") {
		node.CpuName = "INTEL"
	} else if strings.Contains(strings.ToLower(infoStats[0].ModelName), "amd") {
		node.CpuName = "AMD"
	}
	return nil
}
