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
	nodeInfo := new(device.NodeInfo)
	err := device.GetGpu(nodeInfo)
	if err != nil {
		fmt.Printf("The node collect gpu info failed, if the node does not have a GPU, this error can be ignored. %v\n", err)
	}

	err = device.GetHardwareData(nodeInfo)
	if err != nil {
		fmt.Printf("ERROR:: get hardware failed, %v\n", err)
	}

	marshal, err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Printf("convert to json failed, %v \n", err)
	}
	println(string(marshal))
}
