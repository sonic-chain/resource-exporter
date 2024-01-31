package main

import (
	"encoding/json"
	"fmt"
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
	gpuInfo := new(device.NodeInfo)
	err := device.GetGpu(gpuInfo)
	if err != nil {
		fmt.Printf("If the node has a GPU, this error can be ignored. %v \n", err)
		return
	}

	marshal, err := json.Marshal(gpuInfo)
	if err != nil {
		fmt.Printf("ERROR:: convert to json failed, %v \n", err)
		return
	}
	println(string(marshal))
}
