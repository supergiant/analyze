package sunsetting

// CheckAllPodsAtATime makes simple check that it is possible to move all pods of a node to another node.
func CheckAllPodsAtATime(entriesByWastedRam EntriesByWastedRAM) []*InstanceEntry {
	var res = make([]*InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		for i := len(entriesByWastedRam) - 1; i > 0; i-- {
			// check that all requested memory of instance can be moved to another instance
			var wastedRam = entriesByWastedRam[i].AllocatableMemory - entriesByWastedRam[i].MemoryReqs
			if maxWastedRamEntry.MemoryReqs <= wastedRam {
				//sunset this instance
				res = append(res, maxWastedRamEntry)
				//change memory requests of node which receive all workload
				entriesByWastedRam[i].MemoryReqs = entriesByWastedRam[i].MemoryReqs + maxWastedRamEntry.MemoryReqs
				break
			}
		}
	}

	return res
}

func CheckEachPodOneByOne(entriesByWastedRam EntriesByWastedRAM, entriesByRequestedRAM EntriesByRequestedRAM) []*InstanceEntry {
	var res = make([]*InstanceEntry, 0)

	for _, maxWastedRamEntry := range entriesByWastedRam {
		for _, maxRequestedRamEntry := range entriesByRequestedRAM {
			maxWastedRamEntry.RAMRequested()
			maxRequestedRamEntry.RAMWasted()
		}
	}

	return res
}
