package kube

import (
	"strings"

	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/supergiant/robot/builtin/plugins/underutilizednodes/models"
)

func GetInstanceEntries(сoreV1Client *corev1.CoreV1Client) (map[string]*models.InstanceEntry, error) {
	var instanceEntries = map[string]*models.InstanceEntry{}

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

		var entry = &models.InstanceEntry{
			AwsInstance: &models.AwsInstance{},
			KubeWorker: &models.KubeWorker{
				Name: node.Name,
			},
		}

		entry.Region, entry.InstanceID = parseProviderID(node.Spec.ProviderID)

		// calculate minions requests/limits
		reqs, limits := getPodsTotalRequestsAndLimits(nonTerminatedPodsList)
		cpuReqs, cpuLimits := reqs[corev1api.ResourceCPU], limits[corev1api.ResourceCPU]
		memoryReqs, memoryLimits := reqs[corev1api.ResourceMemory], limits[corev1api.ResourceMemory]

		entry.CpuReqs, entry.CpuLimits = cpuReqs.MilliValue(), cpuLimits.MilliValue()
		entry.MemoryReqs, entry.MemoryLimits = memoryReqs.Value(), memoryLimits.Value()

		var allocatable = node.Status.Capacity
		if len(node.Status.Allocatable) > 0 {
			allocatable = node.Status.Allocatable
		}

		entry.AllocatableCpu = allocatable.Cpu().MilliValue()
		entry.AllocatableMemory = allocatable.Memory().Value()

		if entry.AllocatableCpu != 0 {
			entry.FractionCpuReqs = float64(entry.CpuReqs) / float64(entry.AllocatableCpu) * 100
			entry.FractionCpuLimits = float64(entry.CpuLimits) / float64(entry.AllocatableCpu) * 100
		}

		if entry.AllocatableMemory != 0 {
			entry.FractionMemoryReqs = float64(entry.MemoryReqs) / float64(entry.AllocatableMemory) * 100
			entry.FractionMemoryLimits = float64(entry.MemoryLimits) / float64(entry.AllocatableMemory) * 100
		}

		instanceEntries[entry.InstanceID] = entry
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
