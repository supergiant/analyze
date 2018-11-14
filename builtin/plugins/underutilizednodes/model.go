package underutilizednodes

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
	Name string
	AllocatableCpu int64
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

func (m *InstanceEntry) RAMWastedGiB() int64 {
	return m.AllocatableMemory
}

func (m *InstanceEntry) RAMRequestedGiB() int64 {
	return m.MemoryReqs
}

func (m *InstanceEntry) CPURequested() float64 {
	panic("implement me")
}

func (m *InstanceEntry) CPURequestedPercents() float64 {
	panic("implement me")
}

// EntriesByWastedRAM implements sort.Interface based on the KubeWorker.allocatableMemory field.
type EntriesByWastedRAM []*InstanceEntry

func (e EntriesByWastedRAM) Len() int           { return len(e) }
func (e EntriesByWastedRAM) Less(i, j int) bool { return e[i].RAMWastedGiB() < e[j].RAMWastedGiB() }
func (e EntriesByWastedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
