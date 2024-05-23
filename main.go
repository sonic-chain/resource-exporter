package main

import (
	"encoding/json"
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/lagrangedao/resource-exporter/device"
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
	ret := device.GetGpu(nodeInfo)
	if ret != nvml.SUCCESS {
		fmt.Printf("The node initialize nvm libnvidia failed, if the node does not have a GPU, this error can be ignored. code: %d\n", ret)
	}

	err := device.GetHardwareData(nodeInfo)
	if err != nil {
		fmt.Printf("ERROR:: get hardware failed, %v\n", err)
	}

	marshal, err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Printf("ERROR:: convert to json failed, %v \n", err)
	}
	println(string(marshal))
}
