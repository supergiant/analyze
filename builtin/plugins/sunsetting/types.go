package sunsetting

import (
	"sort"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/cloudprovider"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/kube"
)

// InstanceEntry struct represents Kelly's "instances to sunset" table entry,
// It consists of some k8s cluster worker node params ans some ec2 instance params
type InstanceEntry struct {
	CloudProvider cloudprovider.ComputeInstance `json:"cloudProvider"`
	Price         cloudprovider.ProductPrice    `json:"price"`
	WorkerNode    kube.NodeResourceRequirements `json:"kube"`
}

func (m *InstanceEntry) RAMWasted() int64 {
	return m.WorkerNode.AllocatableMemory - m.WorkerNode.MemoryReqs()
}

func (m *InstanceEntry) RAMRequested() int64 {
	return m.WorkerNode.MemoryReqs()
}

func (m *InstanceEntry) CPUWasted() int64 {
	return m.WorkerNode.AllocatableCpu - m.WorkerNode.CpuReqs()
}

// EntriesByWastedRAM implements sort.Interface based on the value returned by NodeResourceRequirements.RAMWasted().
type EntriesByWastedRAM []*InstanceEntry

func NewSortedEntriesByWastedRAM(in []*InstanceEntry) EntriesByWastedRAM {
	var res = make([]*InstanceEntry, len(in))
	for i, e := range in {
		var item = &InstanceEntry{
			CloudProvider: e.CloudProvider,
			Price:         e.Price,
			WorkerNode:    e.WorkerNode,
		}
		item.WorkerNode.PodsResourceRequirements = make([]*kube.PodResourceRequirements, len(e.WorkerNode.PodsResourceRequirements))
		for j, p := range e.WorkerNode.PodsResourceRequirements {
			var newP = *p
			item.WorkerNode.PodsResourceRequirements[j] = &newP
		}
		// we need to calculate totals because sorting logic is based on them
		item.WorkerNode.RefreshTotals()

		res[i] = item
	}
	var entries = EntriesByWastedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByWastedRAM) Len() int           { return len(e) }
func (e EntriesByWastedRAM) Less(i, j int) bool { return e[i].RAMWasted() < e[j].RAMWasted() }
func (e EntriesByWastedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// EntriesByRequestedRAM implements sort.Interface based on the value returned by NodeResourceRequirements.RAMRequested().
type EntriesByRequestedRAM []*InstanceEntry

func NewSortedEntriesByRequestedRAM(in []*InstanceEntry) EntriesByRequestedRAM {
	var res = make([]*InstanceEntry, len(in))
	for i, e := range in {
		var item = &InstanceEntry{
			CloudProvider: e.CloudProvider,
			Price:         e.Price,
			WorkerNode:    e.WorkerNode,
		}
		item.WorkerNode.PodsResourceRequirements = make([]*kube.PodResourceRequirements, len(e.WorkerNode.PodsResourceRequirements))
		for j, p := range e.WorkerNode.PodsResourceRequirements {
			var newP = *p
			item.WorkerNode.PodsResourceRequirements[j] = &newP
		}
		// we need to calculate totals because sorting logic is based on them
		item.WorkerNode.RefreshTotals()
		res[i] = item
	}
	var entries = EntriesByRequestedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByRequestedRAM) Len() int           { return len(e) }
func (e EntriesByRequestedRAM) Less(i, j int) bool { return e[i].RAMRequested() < e[j].RAMRequested() }
func (e EntriesByRequestedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
