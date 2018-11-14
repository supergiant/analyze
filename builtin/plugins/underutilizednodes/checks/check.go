package checks

import "github.com/supergiant/robot/builtin/plugins/underutilizednodes/models"

func AllPodsAtATime(entriesByWastedRam models.EntriesByWastedRAM) []*models.InstanceEntry {
	var res = make([]*models.InstanceEntry, 0)

	for _, maxWatedRamEntry := range entriesByWastedRam {
		for i := len(entriesByWastedRam) - 1; i > 0; i-- {
			// check that all requested memory of instance can be moved to another instance
			var wastedRam = entriesByWastedRam[i].AllocatableMemory - entriesByWastedRam[i].MemoryReqs
			if maxWatedRamEntry.MemoryReqs <= wastedRam {
				//sunset this instance
				res = append(res, maxWatedRamEntry)
				//change memory requests of node which receive all workload
				entriesByWastedRam[i].MemoryReqs = entriesByWastedRam[i].MemoryReqs + maxWatedRamEntry.MemoryReqs
				break
			}
		}
	}

	return res
}
