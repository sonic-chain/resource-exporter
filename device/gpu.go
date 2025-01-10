package device

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func GetGpu(gpu *NodeInfo) error {
	gInfo, err := getGpuInfo()
	if err != nil {
		return err
	}
	if gInfo != nil {
		gpu.Gpu.AttachedGpus = len(gInfo)
		gpu.Gpu.DriverVersion = gInfo[0].driverVersion
	}

	cudaVersion, err := getCudaVersion()
	if err != nil {
		return err
	}
	gpu.Gpu.CudaVersion = cudaVersion

	useGpus, err := getUseGpu()
	if err != nil {
		return err
	}
	for _, gpuU := range useGpus {
		for _, info := range gInfo {
			if gpuU.guid == info.gpuUid {
				info.status = "occupied"
			}
		}
	}

	var gpuDetail []GpuDetail
	for _, info := range gInfo {
		gpuDetail = append(gpuDetail, GpuDetail{
			Index:        info.index,
			OriginalName: info.gpuName,
			ProductName:  strings.ToUpper(convertName(info.gpuName)),
			FbMemoryUsage: Common{
				Total: info.memTotal,
				Used:  info.memUsed,
				Free:  info.memFree,
			},
			Status: info.status,
			Guid:   info.gpuUid,
		})
	}
	gpu.Gpu.Details = gpuDetail
	return nil
}

func getGpuInfo() ([]*collectGpu, error) {
	var cg []*collectGpu
	var err error
	cmd := exec.Command("nvidia-smi", "--query-gpu=index,driver_version,gpu_uuid,gpu_name,memory.total,memory.used,memory.free", "--format=csv,noheader")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed execute nvidia-smi, error:%v", err)
	}

	if len(output) > 0 {
		reader := bufio.NewReader(strings.NewReader(string(output)))
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				break
			}
			fields := strings.Split(string(line), ",")
			cg = append(cg, &collectGpu{
				index:         strings.TrimSpace(fields[0]),
				driverVersion: strings.TrimSpace(fields[1]),
				gpuUid:        strings.TrimSpace(fields[2]),
				gpuName:       strings.TrimSpace(fields[3]),
				memTotal:      strings.TrimSpace(fields[4]),
				memUsed:       strings.TrimSpace(fields[5]),
				memFree:       strings.TrimSpace(fields[6]),
				status:        "available",
			})
		}
	}
	return cg, nil
}

func getUseGpu() ([]usedGpu, error) {
	var ug []usedGpu
	cmd := exec.Command("nvidia-smi", "--query-compute-apps=gpu_uuid,gpu_name,process_name", "--format=csv,noheader")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ug, fmt.Errorf("failed execute nvidia-smi, error:%v", err)
	}

	if len(output) > 0 {
		reader := bufio.NewReader(strings.NewReader(string(output)))
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				break
			}
			fields := strings.Split(string(line), ",")
			ug = append(ug, usedGpu{
				guid:        strings.TrimSpace(fields[0]),
				processName: strings.TrimSpace(fields[2]),
			})
		}
	}
	return ug, nil
}

func getCudaVersion() (string, error) {
	var cudaVersion string

	r, w := io.Pipe()
	cmd := exec.Command("nvidia-smi", "-q")
	cmd.Stdout = w

	grepCmd := exec.Command("grep", "CUDA Version")
	grepCmd.Stdin = r

	var out bytes.Buffer
	grepCmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed starting nvidia-smi command: %v", err)
	}

	err = grepCmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed starting grep command: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("failed waiting for nvidia-smi command: %v", err)
	}
	w.Close()
	err = grepCmd.Wait()
	if err != nil {
		return "", fmt.Errorf("failed waiting for grep command: %v", err)
	}

	if len(out.String()) > 0 && strings.Contains(out.String(), ":") {
		cudaVersion = strings.TrimSpace(strings.Split(out.String(), ":")[1])
	}
	return cudaVersion, nil
}

type collectGpu struct {
	index         string
	driverVersion string
	gpuUid        string
	gpuName       string
	memTotal      string
	memUsed       string
	memFree       string
	status        string //  occupied  available
}

type usedGpu struct {
	guid        string
	processName string
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
