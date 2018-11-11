package underutilizednodes

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/golang/protobuf/ptypes/any"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type uuNodesPlugin struct {
	config *proto.PluginConfig
}

func NewPlugin() proto.PluginClient {
	return &uuNodesPlugin{}
}

func (u *uuNodesPlugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
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

	var minions = []*Minion{}
	var result = &proto.CheckResult{
		ExecutionStatus: "OK",
		Status:          proto.CheckStatus_UNKNOWN_CHECK_STATUS,
		Name:            "Underutilized nodes sunsetting Check",
		Description: &any.Any{
			TypeUrl: "io.supergiant.analyze.plugin.requestslimitscheck",
			Value:   []byte("Resources (CPU/RAM) total capacity and allocatable where checked on nodes of k8s cluster. Results:"),
		},
		Actions: []*proto.Action{
			&proto.Action{
				ActionId:    "1",
				Description: "Dismiss notification",
			},
			&proto.Action{
				ActionId:    "2",
				Description: "Sunset nodes",
			},
		},
	}

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

		var minion = &Minion{
			NonTerminatedPods: nonTerminatedPodsList.Items,
			Node:              node,
		}

		// calculate minions requests/limits
		reqs, limits := getPodsTotalRequestsAndLimits(nonTerminatedPodsList)
		minion.cpuReqs, minion.cpuLimits = reqs[corev1api.ResourceCPU], limits[corev1api.ResourceCPU]
		minion.memoryReqs, minion.memoryLimits = reqs[corev1api.ResourceMemory], limits[corev1api.ResourceMemory]

		var allocatable = node.Status.Capacity
		if len(node.Status.Allocatable) > 0 {
			allocatable = node.Status.Allocatable
		}

		if allocatable.Cpu().MilliValue() != 0 {
			minion.fractionCpuReqs = float64(minion.cpuReqs.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
			minion.fractionCpuLimits = float64(minion.cpuLimits.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
		}

		if allocatable.Memory().Value() != 0 {
			minion.fractionMemoryReqs = float64(minion.memoryReqs.Value()) / float64(allocatable.Memory().Value()) * 100
			minion.fractionMemoryLimits = float64(minion.memoryLimits.Value()) / float64(allocatable.Memory().Value()) * 100
		}

		minions = append(minions, minion)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	var awsConfig = u.config.GetAwsConfig()

	cfg.Credentials = aws.NewStaticCredentialsProvider(
		awsConfig.GetAccessKeyId(),
		awsConfig.GetSecretAccessKey(),
		"",
	)
	cfg.Region = awsConfig.GetRegion()

	svc := ec2.New(cfg)
	var nameFilter = "dmts-tst"

	fmt.Printf("listing instances with tag %v in: %v\n", nameFilter, cfg.Region)
	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("tag:KubernetesCluster"),
				Values: []string{
					strings.Join([]string{"*", nameFilter, "*"}, ""),
				},
			},
		},
	}

	req := svc.DescribeInstancesRequest(params)

	resp, err := req.Send()
	if err != nil {
		fmt.Printf("failed to describe instances, %s, %v", cfg.Region, err)
	}

	if resp != nil {
		for _, r := range resp.Reservations {
			for _, i := range r.Instances {
				fmt.Printf("InstanceId: %v\n", *i.InstanceId)
				fmt.Printf("State: %v\n", *i.State)
				fmt.Printf("InstanceType: %v\n", i.InstanceType)
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						fmt.Printf("Name tag: %v\n", *t.Value)
					}
					if *t.Key == "KubernetesCluster" {
						fmt.Printf("KubernetesCluster tag: %v\n", *t.Value)
					}
				}
			}
		}
	}

	//for _, minion := range minions {
	//	for _, pod := range minion.Pods {
	//		fmt.Printf(
	//			"Pods: Namespace %s, Name %s, fisrtContainerResources %v, nodeName %s \n",
	//			pod.Namespace,
	//			pod.Name,
	//			pod.Spec.Containers[0].Resources,
	//			minion.Node.Name,
	//		)
	//	}
	//}

	return &proto.CheckResponse{Result: result}, nil
}

func (u *uuNodesPlugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *uuNodesPlugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	u.config = in
	//TODO: add here config validation in future
	return nil, nil
}

func (u *uuNodesPlugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *uuNodesPlugin) Info(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
	return &proto.PluginInfo{
		Id:          "supergiant-underutilized-nodes-plugin",
		Version:     "v0.0.1",
		Name:        "Underutilized nodes sunsetting plugin",
		Description: "This plugin checks nodes using high intelligent Kelly's approach to find underutilized nodes, than calculates how it is possible to fix that",
	}, nil
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
