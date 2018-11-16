package kube

import (
	"strings"

	"k8s.io/client-go/rest"

	corev1api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	сoreV1Client *corev1.CoreV1Client
}

func NewKubeClient() (*Client, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the client
	сoreV1Client, err := corev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		сoreV1Client: сoreV1Client,
	}, nil
}

func (c *Client) GetNodeResourceRequirements() (map[string]*NodeResourceRequirements, error) {
	var instanceEntries = map[string]*NodeResourceRequirements{}

	nodes, err := c.сoreV1Client.Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name + ",status.phase!=" + string(corev1api.PodSucceeded) + ",status.phase!=" + string(corev1api.PodFailed))
		if err != nil {
			return nil, err
		}

		nonTerminatedPodsList, err := c.сoreV1Client.Pods("").List(metav1.ListOptions{FieldSelector: fieldSelector.String()})
		if err != nil {
			return nil, err
		}

		var nodeResourceRequirements = getNodeResourceRequirements(node, nonTerminatedPodsList.Items)

		instanceEntries[nodeResourceRequirements.InstanceID] = nodeResourceRequirements
	}

	return instanceEntries, nil
}

func getNodeResourceRequirements(node corev1api.Node, pods []corev1api.Pod) *NodeResourceRequirements {
	var nodeResourceRequirements = &NodeResourceRequirements{
		Name:                     node.Name,
		PodsResourceRequirements: []*PodResourceRequirements{},
	}

	nodeResourceRequirements.Region, nodeResourceRequirements.InstanceID = parseProviderID(node.Spec.ProviderID)

	// calculate worker node requests/limits
	nodeResourceRequirements.PodsResourceRequirements = getPodsRequestsAndLimits(pods)

	var allocatable = node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}

	nodeResourceRequirements.AllocatableCpu = allocatable.Cpu().MilliValue()
	nodeResourceRequirements.AllocatableMemory = allocatable.Memory().Value()

	nodeResourceRequirements.RefreshTotals()

	return nodeResourceRequirements
}

// TODO: add checks and errors
// for aws ProviderID has format - aws:///us-west-1b/i-0c912bfd4048b97e5
func parseProviderID(providerID string) (string, string) {
	var s = strings.TrimPrefix(providerID, "aws:///")
	ss := strings.Split(s, "/")
	return ss[0], ss[1]
}

func getPodsRequestsAndLimits(podList []corev1api.Pod) []*PodResourceRequirements {
	var result = []*PodResourceRequirements{}
	for _, pod := range podList {
		var podRR = &PodResourceRequirements{
			PodName: pod.Name,
		}

		podReqs, podLimits := PodRequestsAndLimits(&pod)
		cpuReqs, cpuLimits := podReqs[corev1api.ResourceCPU], podLimits[corev1api.ResourceCPU]
		memoryReqs, memoryLimits := podReqs[corev1api.ResourceMemory], podLimits[corev1api.ResourceMemory]
		podRR.CpuReqs, podRR.CpuLimits = cpuReqs.MilliValue(), cpuLimits.MilliValue()
		podRR.MemoryReqs, podRR.MemoryLimits = memoryReqs.Value(), memoryLimits.Value()

		result = append(result, podRR)
	}

	return result
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
