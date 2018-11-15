package kube

import (
	"strings"

	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func GetNodeResourceRequirements(сoreV1Client *corev1.CoreV1Client) (map[string]*NodeResourceRequirements, error) {
	var instanceEntries = map[string]*NodeResourceRequirements{}

	nodes, err := сoreV1Client.Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name + ",status.phase!=" + string(corev1api.PodSucceeded) + ",status.phase!=" + string(corev1api.PodFailed))
		if err != nil {
			return nil, err
		}

		nonTerminatedPodsList, err := сoreV1Client.Pods("").List(metav1.ListOptions{FieldSelector: fieldSelector.String()})
		if err != nil {
			return nil, err
		}

		var nodeResourceRequirements = &NodeResourceRequirements{
			Name: node.Name,
		}

		nodeResourceRequirements.Region, nodeResourceRequirements.InstanceID = parseProviderID(node.Spec.ProviderID)

		// calculate minions requests/limits
		reqs, limits := getPodsTotalRequestsAndLimits(nonTerminatedPodsList)
		cpuReqs, cpuLimits := reqs[corev1api.ResourceCPU], limits[corev1api.ResourceCPU]
		memoryReqs, memoryLimits := reqs[corev1api.ResourceMemory], limits[corev1api.ResourceMemory]

		nodeResourceRequirements.CpuReqs, nodeResourceRequirements.CpuLimits = cpuReqs.MilliValue(), cpuLimits.MilliValue()
		nodeResourceRequirements.MemoryReqs, nodeResourceRequirements.MemoryLimits = memoryReqs.Value(), memoryLimits.Value()

		var allocatable = node.Status.Capacity
		if len(node.Status.Allocatable) > 0 {
			allocatable = node.Status.Allocatable
		}

		nodeResourceRequirements.AllocatableCpu = allocatable.Cpu().MilliValue()
		nodeResourceRequirements.AllocatableMemory = allocatable.Memory().Value()

		if nodeResourceRequirements.AllocatableCpu != 0 {
			nodeResourceRequirements.FractionCpuReqs = float64(nodeResourceRequirements.CpuReqs) / float64(nodeResourceRequirements.AllocatableCpu) * 100
			nodeResourceRequirements.FractionCpuLimits = float64(nodeResourceRequirements.CpuLimits) / float64(nodeResourceRequirements.AllocatableCpu) * 100
		}

		if nodeResourceRequirements.AllocatableMemory != 0 {
			nodeResourceRequirements.FractionMemoryReqs = float64(nodeResourceRequirements.MemoryReqs) / float64(nodeResourceRequirements.AllocatableMemory) * 100
			nodeResourceRequirements.FractionMemoryLimits = float64(nodeResourceRequirements.MemoryLimits) / float64(nodeResourceRequirements.AllocatableMemory) * 100
		}

		instanceEntries[nodeResourceRequirements.InstanceID] = nodeResourceRequirements
	}

	return instanceEntries, nil
}

// TODO: add checks and errors
// for aws ProviderID has format - aws:///us-west-1b/i-0c912bfd4048b97e5
func parseProviderID(providerID string) (string, string) {
	var s = strings.TrimPrefix(providerID, "aws:///")
	ss := strings.Split(s, "/")
	return ss[0], ss[1]
}

func getPodsTotalRequestsAndLimits(podList *corev1api.PodList) (reqs map[corev1api.ResourceName]resource.Quantity, limits map[corev1api.ResourceName]resource.Quantity) {
	reqs, limits = map[corev1api.ResourceName]resource.Quantity{}, map[corev1api.ResourceName]resource.Quantity{}
	for _, pod := range podList.Items {
		podReqs, podLimits := PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = *podReqValue.Copy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = *podLimitValue.Copy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	return
}

// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod.
func PodRequestsAndLimits(pod *corev1api.Pod) (reqs corev1api.ResourceList, limits corev1api.ResourceList) {
	reqs, limits = corev1api.ResourceList{}, corev1api.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}
	return
}

// addResourceList adds the resources in newList to list
func addResourceList(list, new corev1api.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, new corev1api.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = *quantity.Copy()
			}
		}
	}
}
