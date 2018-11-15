package sunsetting

import (
	"sort"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/prices"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/kube"
)

// InstanceEntry struct represents Kelly's "instances to sunset" table entry,
// It consists of some k8s cluster worker node params ans some ec2 instance params
type InstanceEntry struct {
	//AWS instance Type
	InstanceType string
	Price        *prices.Item
	*kube.NodeResourceRequirements
}

func (m *InstanceEntry) RAMWasted() int64 {
	return m.AllocatableMemory - m.MemoryReqs()
}

func (m *InstanceEntry) RAMRequested() int64 {
	return m.MemoryReqs()
}

func (m *InstanceEntry) CPURequested() float64 {
	panic("implement me")
}

func (m *InstanceEntry) CPURequestedPercents() float64 {
	panic("implement me")
}

// EntriesByWastedRAM implements sort.Interface based on the value returned by NodeResourceRequirements.RAMWasted().
type EntriesByWastedRAM []*InstanceEntry

func NewSortedEntriesByWastedRAM(in []*InstanceEntry) EntriesByWastedRAM {
	var res = make([]*InstanceEntry, len(in))
	for i, e := range in {
		var item = &InstanceEntry{
			InstanceType:             e.InstanceType,
			Price:                    &prices.Item{},
			NodeResourceRequirements: &kube.NodeResourceRequirements{},
		}
		*item.Price = *e.Price
		*item.NodeResourceRequirements = *e.NodeResourceRequirements
		for _, p := range e.NodeResourceRequirements.PodsResourceRequirements {
			var newP = *p
			item.NodeResourceRequirements.PodsResourceRequirements = append(item.NodeResourceRequirements.PodsResourceRequirements, &newP)
		}

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
			InstanceType:             e.InstanceType,
			Price:                    &prices.Item{},
			NodeResourceRequirements: &kube.NodeResourceRequirements{},
		}
		*item.Price = *e.Price
		*item.NodeResourceRequirements = *e.NodeResourceRequirements
		for _, p := range e.NodeResourceRequirements.PodsResourceRequirements {
			var newP = *p
			item.NodeResourceRequirements.PodsResourceRequirements = append(item.NodeResourceRequirements.PodsResourceRequirements, &newP)
		}

		res[i] = item
	}
	var entries = EntriesByRequestedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByRequestedRAM) Len() int           { return len(e) }
func (e EntriesByRequestedRAM) Less(i, j int) bool { return e[i].RAMRequested() < e[j].RAMRequested() }
func (e EntriesByRequestedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
