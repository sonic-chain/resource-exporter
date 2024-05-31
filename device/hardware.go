package device

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetHardwareData(node *NodeInfo) error {

	diskStat, err := getDiskUsage("/")
	if err != nil {
		return err
	}
	node.Storage.Total = diskStat.Total
	node.Storage.Used = diskStat.Used
	node.Storage.Free = diskStat.Available

	mem, err := getMemoryUsage()
	if err != nil {
		return err
	}
	node.Memory.Total = mem.Total
	node.Memory.Used = mem.Used
	node.Memory.Free = mem.Available

	totalCores, totalUsage, availableUsage, err := getCpuUsage()
	if err != nil {
		return err
	}
	node.Cpu.Total = fmt.Sprintf("%d", totalCores)
	node.Cpu.Used = fmt.Sprintf("%d", totalUsage)
	node.Cpu.Free = fmt.Sprintf("%d", availableUsage)

	node.Vcpu.Total = node.Cpu.Total
	node.Vcpu.Used = node.Cpu.Used
	node.Vcpu.Free = node.Cpu.Free

	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return err
	}
	uuid := strings.TrimSpace(string(data))
	node.MachineId = fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])

	arch := runtime.GOARCH
	if strings.Contains(strings.ToLower(arch), "intel") {
		node.CpuName = "INTEL"
	} else if strings.Contains(strings.ToLower(arch), "amd") {
		node.CpuName = "AMD"
	}
	return nil
}

func getCpuUsage() (totalCores int, totalUsage int, availableUsage int, err error) {
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	totalCores, err = strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, 0, 0, err
	}

	first, err := checkCpuUsage()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)
	second, err := checkCpuUsage()
	if err != nil {
		log.Fatal(err)
	}

	idle := float64(second.idle - first.idle)
	used := float64(second.used - first.used)
	var usage float64
	if idle+used > 0 {
		usage = used / (idle + used)
	}

	totalUsage = int(math.Round(float64(totalCores) * usage))
	if totalUsage >= totalCores {
		totalUsage = totalCores
	}
	availableUsage = totalCores - totalUsage

	return totalCores, totalUsage, availableUsage, nil
}

type DiskUsage struct {
	Total     string
	Used      string
	Available string
}

func getDiskUsage(mountOn string) (*DiskUsage, error) {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var diskUsages = new(DiskUsage)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		if fields[5] == mountOn {
			diskUsages.Total = fields[1]
			diskUsages.Used = fields[2]
			diskUsages.Available = fields[3]
		}
	}
	return diskUsages, nil
}

type MemoryUsage struct {
	Total     string
	Used      string
	Available string
}

func getMemoryUsage() (*MemoryUsage, error) {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected output from free command")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 7 {
		return nil, fmt.Errorf("unexpected output from free command")
	}
	total, _ := strconv.ParseUint(fields[1], 10, 64)
	used, _ := strconv.ParseUint(fields[2], 10, 64)
	available, _ := strconv.ParseUint(fields[6], 10, 64)

	return &MemoryUsage{
		Total:     fmt.Sprintf("%d GiB", total/1024),
		Used:      fmt.Sprintf("%d GiB", used/1024),
		Available: fmt.Sprintf("%d GiB", available/1024),
	}, nil
}

type result struct {
	used uint64
	idle uint64
}

func checkCpuUsage() (*result, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	res := &result{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "cpu" {
			continue
		}
		n := len(fields)
		for i := 1; i < n; i++ {
			if i > 8 {
				continue
			}

			val, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				return nil, err
			}
			if i == 4 || i == 5 {
				res.idle += val
			} else {
				res.used += val
			}
		}
		return res, nil
	}
	return res, nil
}
