package main

import (
	"encoding/json"
	"fmt"
	"github.com/lagrangedao/resource-exporter/device"
	"os/exec"
	"strings"
	"time"
)

func main() {
	printLog()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			printLog()
		case <-time.After(5 * time.Second):
			break
		}
	}
}

func printLog() {
	nodeInfo := new(device.NodeInfo)
	err := device.GetGpu(nodeInfo)
	if err != nil {
		if strings.Contains(err.Error(), "12") {
			fmt.Printf("The node not found nvm libnvidia, if the node does not have a GPU, this error can be ignored.\n")
		} else {
			fmt.Printf("%+v\n", err)
			return
		}
	}

	cmd := exec.Command("sh", "-c", "lscpu | grep 'Model name'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR:: executing command: %v\n", err)
		return
	}
	if strings.Contains(strings.ToUpper(string(output)), "AMD") {
		nodeInfo.CpuName = "AMD"
	} else if strings.Contains(strings.ToUpper(string(output)), "INTEL") {
		nodeInfo.CpuName = "INTEL"
	}

	marshal, err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Printf("ERROR:: convert to json failed, %v \n", err)
		return
	}
	println(string(marshal))
}
