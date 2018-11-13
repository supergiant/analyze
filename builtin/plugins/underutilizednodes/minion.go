package underutilizednodes

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Minion struct represents k8s cluster worker node
type Minion struct {
	*AwsInstance
	*KubeWorker
}

type AwsInstance struct {
	Region       string
	InstanceID   string
	InstanceType string
}

type KubeWorker struct {
	NonTerminatedPods []v1.Pod
	Node              *v1.Node

	cpuReqs      resource.Quantity
	cpuLimits    resource.Quantity
	memoryReqs   resource.Quantity
	memoryLimits resource.Quantity

	fractionCpuReqs      float64
	fractionCpuLimits    float64
	fractionMemoryReqs   float64
	fractionMemoryLimits float64
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
