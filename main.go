package main

import (
	"encoding/json"
	"github.com/lagrangedao/resource-exporter/device"
	"log"
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
	device.GetGpu(gpuInfo)

	marshal, err := json.Marshal(gpuInfo)
	if err != nil {
		log.Println(err)
		return
	}
	println(string(marshal))
}
