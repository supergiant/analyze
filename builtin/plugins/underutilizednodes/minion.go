package underutilizednodes

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Minion struct represents k8s cluster worker node
type Minion struct {
	AWSZone           string
	InstanceID        string
	InstanceType      string
	NonTerminatedPods []v1.Pod
	Node              v1.Node

	cpuReqs, cpuLimits, memoryReqs, memoryLimits                                 resource.Quantity
	fractionCpuReqs, fractionCpuLimits, fractionMemoryReqs, fractionMemoryLimits float64
}

func (m *Minion) RAMRequestedGiB() float64 {
	panic("implement me")
}

func (m *Minion) RAMRequestedPercents() float64 {
	panic("implement me")
}

func (m *Minion) CPURequested() float64 {
	panic("implement me")
}

func (m *Minion) CPURequestedPercents() float64 {
	panic("implement me")
}
