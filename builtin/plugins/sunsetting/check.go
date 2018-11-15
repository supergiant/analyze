package sunsetting

import "sort"

// CheckAllPodsAtATime makes simple check that it is possible to move all pods of a node to another node.
func CheckAllPodsAtATime(entriesByWastedRam EntriesByWastedRAM) []*InstanceEntry {
	var res = make([]*InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		for i := len(entriesByWastedRam) - 1; i > 0; i-- {
			// check that all requested memory of instance can be moved to another instance
			var wastedRam = entriesByWastedRam[i].AllocatableMemory - entriesByWastedRam[i].MemoryReqs()
			if maxWastedRamEntry.MemoryReqs() <= wastedRam {
				//sunset this instance
				res = append(res, maxWastedRamEntry)
				//change memory requests of node which receive all workload
				entriesByWastedRam[i].PodsResourceRequirements = append(entriesByWastedRam[i].PodsResourceRequirements, maxWastedRamEntry.PodsResourceRequirements...)
				entriesByWastedRam[i].RefreshTotals()
				break
			}
		}
	}

	return res
}

func CheckEachPodOneByOne(entriesByWastedRam EntriesByWastedRAM, entriesByRequestedRAM EntriesByRequestedRAM) []*InstanceEntry {
	var res = make([]*InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		// sort pods in descending order by requested memory
		sort.Slice(
			maxWastedRamEntry.PodsResourceRequirements,
			func(i, j int) bool {
				return maxWastedRamEntry.PodsResourceRequirements[i].MemoryReqs > maxWastedRamEntry.PodsResourceRequirements[j].MemoryReqs
			},
		)

		// check
		for i := 0; i < len(maxWastedRamEntry.PodsResourceRequirements); i++ {
			var podRR = maxWastedRamEntry.PodsResourceRequirements[i]
			for _, maxRequestedRamEntry := range entriesByRequestedRAM {
				if maxRequestedRamEntry.RAMWasted() >= podRR.MemoryReqs {

				}
			}
		}
	}

	return res
}
