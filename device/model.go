package device

type NodeInfo struct {
	Gpu       Gpu    `json:"gpu"`
	MachineId string `json:"machine_id"`
	CpuName   string `json:"cpu_name"`
	Cpu       Common `json:"cpu"`
	Vcpu      Common `json:"vcpu"`
	Memory    Common `json:"memory"`
	Storage   Common `json:"storage"`
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
	Status          string `json:"status"`
}

type Common struct {
	Total string `json:"total"`
	Used  string `json:"used"`
	Free  string `json:"free"`
}
