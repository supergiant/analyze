package sunsetting

import "sort"

// CheckAllPodsAtATime makes simple check that it is possible to move all pods of a node to another node.
func CheckAllPodsAtATime(unsortedEntries []*InstanceEntry) []InstanceEntry {
	var entriesByWastedRam = NewSortedEntriesByWastedRAM(unsortedEntries)
	var res = make([]InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		maxWastedRamEntry.WorkerNode.RefreshTotals()

		for i := len(entriesByWastedRam) - 1; i > 0; i-- {
			var instanceWithMinimumWastedRam = entriesByWastedRam[i]
			instanceWithMinimumWastedRam.WorkerNode.RefreshTotals()
			// check that all requested memory of instance can be moved to another instance
			var wastedRam = instanceWithMinimumWastedRam.RAMWasted()
			var wastedCpu = instanceWithMinimumWastedRam.CPUWasted()
			if maxWastedRamEntry.WorkerNode.MemoryReqs() <= wastedRam && maxWastedRamEntry.WorkerNode.CpuReqs() <= wastedCpu {
				//sunset this instance
				res = append(res, *maxWastedRamEntry)
				//change memory requests of node which receive all workload
				instanceWithMinimumWastedRam.WorkerNode.PodsResourceRequirements = append(instanceWithMinimumWastedRam.WorkerNode.PodsResourceRequirements, maxWastedRamEntry.WorkerNode.PodsResourceRequirements...)
				instanceWithMinimumWastedRam.WorkerNode.RefreshTotals()

				maxWastedRamEntry.WorkerNode.PodsResourceRequirements = nil
				maxWastedRamEntry.WorkerNode.RefreshTotals()
				break
			}
		}
	}

	return res
}

func CheckEachPodOneByOne(unsortedEntries []*InstanceEntry) []InstanceEntry {
	var entriesByWastedRam = NewSortedEntriesByWastedRAM(unsortedEntries)
	var entriesByRequestedRAM = NewSortedEntriesByRequestedRAM(unsortedEntries)
	var res = make([]InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		// sort pods in descending order by requested memory to more effective pods packing on nodes
		sort.Slice(
			maxWastedRamEntry.WorkerNode.PodsResourceRequirements,
			func(i, j int) bool {
				return maxWastedRamEntry.WorkerNode.PodsResourceRequirements[i].MemoryReqs > maxWastedRamEntry.WorkerNode.PodsResourceRequirements[j].MemoryReqs
			},
		)

		// check
		for i := 0; i < len(maxWastedRamEntry.WorkerNode.PodsResourceRequirements); i++ {
			var podRR = maxWastedRamEntry.WorkerNode.PodsResourceRequirements[i]
			for _, maxRequestedRamEntry := range entriesByRequestedRAM {
				if maxRequestedRamEntry.RAMWasted() >= podRR.MemoryReqs && maxRequestedRamEntry.CPUWasted() >= podRR.CpuReqs {
					// we can move the pod
					// delete it from  maxWastedRamEntry
					maxWastedRamEntry.WorkerNode.PodsResourceRequirements = append(maxWastedRamEntry.WorkerNode.PodsResourceRequirements[:i], maxWastedRamEntry.WorkerNode.PodsResourceRequirements[i+1:]...)
					maxWastedRamEntry.WorkerNode.RefreshTotals()
					//and add to maxRequestedRamEntry
					maxRequestedRamEntry.WorkerNode.PodsResourceRequirements = append(maxRequestedRamEntry.WorkerNode.PodsResourceRequirements, podRR)
					maxRequestedRamEntry.WorkerNode.RefreshTotals()
					// if last pod was moved out exit from loop and basically because of i == o we also exit from outer loop
					if len(maxWastedRamEntry.WorkerNode.PodsResourceRequirements) == 0 {
						break
					}
					// reset index to zero because after maxWastedRamEntry.WorkerNode.PodsResourceRequirements slice modification we need to start from scratch
					i = 0
					continue
				}
			}
		}
		if len(maxWastedRamEntry.WorkerNode.PodsResourceRequirements) == 0 {
			res = append(res, *maxWastedRamEntry)
		}
	}

	return res
}
