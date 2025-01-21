package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lagrangedao/resource-exporter/device"
	"time"
)

func main() {
	printLog()

	ticker := time.NewTicker(30 * time.Second)
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
	nodeInfo.TimeStamp = time.Now().Unix()
	originData := fmt.Sprintf("%s+%s+%s+%s$%d", nodeInfo.MachineId, nodeInfo.CpuName,
		nodeInfo.Gpu.DriverVersion, nodeInfo.Cpu.Total, nodeInfo.TimeStamp)
	nodeInfo.CheckCode = hashSHA256(originData)
	marshal, err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Printf("convert to json failed, %v \n", err)
	}
	println(string(marshal))
}

func hashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
