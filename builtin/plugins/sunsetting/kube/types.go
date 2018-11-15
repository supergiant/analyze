package kube

type NodeResourceRequirements struct {
	Name string
	// extracted from ProviderID
	Region     string
	InstanceID string

	Pods              []*PodResourceRequirements
	AllocatableCpu    int64
	AllocatableMemory int64

	CpuReqs      int64
	CpuLimits    int64
	MemoryReqs   int64
	MemoryLimits int64

	FractionCpuReqs      float64
	FractionCpuLimits    float64
	FractionMemoryReqs   float64
	FractionMemoryLimits float64
}

type PodResourceRequirements struct {
	CpuReqs      int64
	CpuLimits    int64
	MemoryReqs   int64
	MemoryLimits int64
}
