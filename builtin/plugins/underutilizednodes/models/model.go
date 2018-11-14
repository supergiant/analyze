package models

import (
	"sort"
)

type PriceItem struct {
	InstanceType string
	Memory       string
	Vcpu         string
	Unit         string
	Currency     string
	ValuePerUnit string
	UsageType    string
	Tenancy      string
}

// InstanceEntry struct represents Kelly's "instances to sunset" table entry,
// It consists of some k8s cluster worker node params ans some ec2 instance params
type InstanceEntry struct {
	*AwsInstance
	*KubeWorker
}

type AwsInstance struct {
	Region       string
	InstanceID   string
	InstanceType string
}

type KubeWorker struct {
	Name              string
	Pods              []*Pod
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

type Pod struct {
	CpuReqs      int64
	CpuLimits    int64
	MemoryReqs   int64
	MemoryLimits int64
}

func (m *InstanceEntry) RAMWasted() int64 {
	return m.AllocatableMemory - m.MemoryReqs
}

func (m *InstanceEntry) RAMRequested() int64 {
	return m.MemoryReqs
}

func (m *InstanceEntry) CPURequested() float64 {
	panic("implement me")
}

func (m *InstanceEntry) CPURequestedPercents() float64 {
	panic("implement me")
}

// EntriesByWastedRAM implements sort.Interface based on the value returned by KubeWorker.RAMWasted().
type EntriesByWastedRAM []*InstanceEntry

func NewSortedEntriesByWastedRAM(in []*InstanceEntry) EntriesByWastedRAM {
	var res = make([]*InstanceEntry, len(in))
	for i, e := range in {
		res[i] = e
	}
	var entries = EntriesByWastedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByWastedRAM) Len() int           { return len(e) }
func (e EntriesByWastedRAM) Less(i, j int) bool { return e[i].RAMWasted() < e[j].RAMWasted() }
func (e EntriesByWastedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// EntriesByRequestedRAM implements sort.Interface based on the value returned by KubeWorker.RAMRequested().
type EntriesByRequestedRAM []*InstanceEntry

func NewSortedEntriesByRequestedRAM(in []*InstanceEntry) EntriesByRequestedRAM {
	var res = make([]*InstanceEntry, len(in))
	for i, e := range in {
		res[i] = e
	}
	var entries = EntriesByRequestedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByRequestedRAM) Len() int           { return len(e) }
func (e EntriesByRequestedRAM) Less(i, j int) bool { return e[i].RAMRequested() < e[j].RAMRequested() }
func (e EntriesByRequestedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
