package device

type NodeInfo struct {
	//MachineId string `json:"machine_id"`
	//Hostname string `json:"hostname"`
	Gpu Gpu `json:"gpu"`
}

type Gpu struct {
	DriverVersion string      `json:"driver_version"`
	CudaVersion   string      `json:"cuda_version"`
	AttachedGpus  int         `json:"attached_gpus"`
	Details       []GpuDetail `json:"details"`
}

type GpuDetail struct {
	ProductName     string `json:"product_name"`
	FbMemoryUsage   Common `json:"fb_memory_usage"`
	Bar1MemoryUsage Common `json:"bar1_memory_usage"`
}

type Common struct {
	Total string `json:"total"`
	Used  string `json:"used"`
	Free  string `json:"free"`
}
