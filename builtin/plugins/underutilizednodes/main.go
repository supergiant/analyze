package underutilizednodes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/builtin/plugins/underutilizednodes/checks"
	"github.com/supergiant/robot/builtin/plugins/underutilizednodes/models"
	"github.com/supergiant/robot/builtin/plugins/underutilizednodes/prices"
	"github.com/supergiant/robot/pkg/plugin/proto"
)

type plugin struct {
	config     *proto.PluginConfig
	ec2Service *ec2.EC2

	computeInstancesPrices map[string][]models.PriceItem
}

var checkResult = &proto.CheckResult{
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

func NewPlugin() proto.PluginClient {
	return &plugin{}
}

func (u *plugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
	var instanceEntries, err = u.getInstanceEntries()
	if err != nil {
		fmt.Printf("unable to get instanceEntries, %v", err)
		return nil, errors.Wrap(err, "unable to get instanceEntries")
	}

	ec2Reservations, err := u.getEC2Reservations()
	if err != nil {
		fmt.Printf("failed to describe ec2 instances, %v", err)
		return nil, errors.Wrap(err, "failed to describe ec2 instances")
	}

	// enrich instanceEntries with ec2 instance type info
	for _, instancesReservation := range ec2Reservations {
		for _, i := range instancesReservation.Instances {
			if i.InstanceId == nil {
				continue
			}
			var minion, exists = instanceEntries[*i.InstanceId]
			if !exists {
				continue
			}

			var instanceType, _ = i.InstanceType.MarshalValue()
			minion.InstanceType = instanceType
		}
	}

	var unsortedEntires []*models.InstanceEntry

	for _, entry := range instanceEntries {
		unsortedEntires = append(unsortedEntires, entry)
	}

	var sortedByWastedRam = models.NewSortedEntriesByWastedRAM(unsortedEntires)
	//var sortedByRequestedRam = models.NewSortedEntriesByRequestedRAM(unsortedEntires)

	var instancesToSunset = checks.AllPodsAtATime(sortedByWastedRam)

	b, _ := json.Marshal(instancesToSunset)

	checkResult.Description = &any.Any{
		TypeUrl: "test",
		Value:   b,
	}
	checkResult.Status = proto.CheckStatus_GREEN
	return &proto.CheckResponse{Result: checkResult}, nil
}

// TODO: add checks and errors
// for aws ProviderID has format - aws:///us-west-1b/i-0c912bfd4048b97e5
func parseProviderID(providerID string) (string, string) {
	var s = strings.TrimPrefix(providerID, "aws:///")
	ss := strings.Split(s, "/")
	return ss[0], ss[1]
}

func (u *plugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *plugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	u.config = in
	//TODO: add here config validation in future

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("unable to load AWS SDK config,  %v", err)
		return nil, errors.Wrap(err, "unable to load AWS SDK config")
	}

	var awsConfig = u.config.GetAwsConfig()
	cfg.Region = awsConfig.GetRegion()
	// TODO bug in sdk?
	cfg.Region = "us-east-1"
	var pricingService = pricing.New(cfg)

	//TODO may be add some init method to plugin?
	u.computeInstancesPrices = prices.Get(pricingService, cfg.Region)

	cfg.Credentials = aws.NewStaticCredentialsProvider(
		awsConfig.GetAccessKeyId(),
		awsConfig.GetSecretAccessKey(),
		"",
	)
	cfg.Region = awsConfig.GetRegion()
	u.ec2Service = ec2.New(cfg)

	return nil, nil
}

func (u *plugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *plugin) Info(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
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

func (u *plugin) getInstanceEntries() (map[string]*models.InstanceEntry, error) {
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

func (u *plugin) getEC2Reservations() ([]ec2.RunInstancesOutput, error) {

	instancesRequest := u.ec2Service.DescribeInstancesRequest(nil)

	describeInstancesResponse, err := instancesRequest.Send()
	if err != nil {
		return nil, err
	}

	return describeInstancesResponse.Reservations, nil
}
