package kube

import "encoding/json"

type NodeResourceRequirements struct {
	Name string
	// extracted from ProviderID
	Region     string
	InstanceID string

	PodsResourceRequirements []*PodResourceRequirements
	AllocatableCpu           int64
	AllocatableMemory        int64

	cpuReqs      int64
	cpuLimits    int64
	memoryReqs   int64
	memoryLimits int64

	fractionCpuReqs      float64
	fractionCpuLimits    float64
	fractionMemoryReqs   float64
	fractionMemoryLimits float64
}

type PodResourceRequirements struct {
	PodName      string `json:"podName"`
	CpuReqs      int64  `json:"cpuRequests"`
	CpuLimits    int64  `json:"cpuLimits"`
	MemoryReqs   int64  `json:"memoryRequests"`
	MemoryLimits int64  `json:"memoryLimits"`
}

// RefreshTotals recalculates total node requests and limits and their fractional representation
// need to be invoked every time when PodsResourceRequirements or AllocatableMemory or AllocatableCpu where changed
func (n *NodeResourceRequirements) RefreshTotals() {
	for _, podRR := range n.PodsResourceRequirements {
		n.cpuReqs += podRR.CpuReqs
		n.cpuLimits += podRR.CpuLimits
		n.memoryReqs += podRR.MemoryReqs
		n.memoryLimits += podRR.MemoryLimits
	}

	if n.AllocatableCpu != 0 {
		n.fractionCpuReqs = float64(n.cpuReqs) / float64(n.AllocatableCpu) * 100
		n.fractionCpuLimits = float64(n.cpuLimits) / float64(n.AllocatableCpu) * 100
	}

	if n.AllocatableMemory != 0 {
		n.fractionMemoryReqs = float64(n.memoryReqs) / float64(n.AllocatableMemory) * 100
		n.fractionMemoryLimits = float64(n.memoryLimits) / float64(n.AllocatableMemory) * 100
	}
}

func (n *NodeResourceRequirements) CpuReqs() int64 {
	return n.cpuReqs
}

func (n *NodeResourceRequirements) CpuLimits() int64 {
	return n.cpuLimits
}

func (n *NodeResourceRequirements) MemoryReqs() int64 {
	return n.memoryReqs
}

func (n *NodeResourceRequirements) MemoryLimits() int64 {
	return n.memoryLimits
}

func (n *NodeResourceRequirements) FractionCpuReqs() float64 {
	return n.fractionCpuReqs
}

func (n *NodeResourceRequirements) FractionCpuLimits() float64 {
	return n.fractionCpuLimits
}

func (n *NodeResourceRequirements) FractionMemoryReqs() float64 {
	return n.fractionMemoryReqs
}

func (n *NodeResourceRequirements) FractionMemoryLimits() float64 {
	return n.fractionMemoryLimits
}

func (n *NodeResourceRequirements) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":                     n.Name,
		"region":                   n.Region,
		"instanceId":               n.InstanceID,
		"podsResourceRequirements": n.PodsResourceRequirements,
		"allocatableCpu":           n.AllocatableCpu,
		"allocatableMemory":        n.AllocatableMemory,
		"cpuRequests":              n.cpuReqs,
		"cpuLimits":                n.cpuLimits,
		"memoryRequests":           n.memoryReqs,
		"memoryLimits":             n.memoryLimits,
		"fractionCpuRequests":      n.fractionCpuReqs,
		"fractionCpuLimits":        n.fractionCpuLimits,
		"fractionMemoryRequests":   n.fractionMemoryReqs,
		"fractionMemoryLimits":     n.fractionMemoryLimits,
	})
}
