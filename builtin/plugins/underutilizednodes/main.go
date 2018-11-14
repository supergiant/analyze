package underutilizednodes

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	config                 *proto.PluginConfig
	computeInstancesPrices map[string][]priceItem
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
	return &uuNodesPlugin{
		computeInstancesPrices: make(map[string][]priceItem, 0),
	}
}

func (u *uuNodesPlugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
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

	var unsortedEntires []*InstanceEntry

	for _, entry := range instanceEntries {
		unsortedEntires = append(unsortedEntires, entry)
	}

	sortedByWastedRam := EntriesByWastedRAM(unsortedEntires)
	sort.Sort(sortedByWastedRam)

	fmt.Printf("sortedByWastedRam:, %v", sortedByWastedRam)

	b, _ := json.Marshal(sortedByWastedRam)

	checkResult.Description = &any.Any{
		TypeUrl:              "test",
		Value:                b,
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

func (u *uuNodesPlugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *uuNodesPlugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	u.config = in
	//TODO: add here config validation in future

	u.getPrices()

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

// TODO add checks and return error
func getProduct(productItem aws.JSONValue) priceItem {
	var pi = priceItem{}
	productInterface, exists := productItem["product"]
	if !exists {
		fmt.Printf("product elemnt doesn't exist")
		return pi
	}

	product, ok := productInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("product elemnt is not map")
		return pi
	}

	attributes, exists := product["attributes"]
	if !exists {
		fmt.Printf("product elemnt doesn't exist")
		return pi
	}

	attrs, ok := attributes.(map[string]interface{})
	if !ok {
		fmt.Printf("attributes elemnt doesn't exist")
		return pi
	}

	value := attrs["instanceType"]
	pi.instanceType, _ = value.(string)
	value = attrs["memory"]
	pi.memory, _ = value.(string)
	value = attrs["vcpu"]
	pi.vcpu, _ = value.(string)
	value = attrs["usagetype"]
	pi.usageType, _ = value.(string)
	value = attrs["tenancy"]
	pi.tenancy, _ = value.(string)

	termsInterface, exists := productItem["terms"]
	if !exists {
		fmt.Printf("terms elemnt doesn't exist")
		return pi
	}

	terms, ok := termsInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("terms elemnt is not map")
		return pi
	}

	onDemandInterface, exists := terms["OnDemand"]
	if !exists {
		fmt.Printf("OnDemand elemnt doesn't exist")
		return pi
	}

	onDemand, ok := onDemandInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("onDemand elemnt is not map")
		return pi
	}

	for _, skuValueInterface := range onDemand {
		skuValue, ok := skuValueInterface.(map[string]interface{})
		if !ok {
			fmt.Printf("skuValue elemnt is not map")
			return pi
		}

		priceDimensionsInterface, exists := skuValue["priceDimensions"]
		if !exists {
			fmt.Printf("priceDimensions elemnt doesn't exist")
			return pi
		}

		priceDimensions, ok := priceDimensionsInterface.(map[string]interface{})
		if !ok {
			fmt.Printf("priceDimensions elemnt is not map")
			return pi
		}

		for _, priceDimentionInterface := range priceDimensions {
			priceDimention, ok := priceDimentionInterface.(map[string]interface{})
			if !ok {
				fmt.Printf("priceDimention elemnt is not map")
				return pi
			}

			unitInterface, exists := priceDimention["unit"]
			if !exists {
				fmt.Printf("unit elemnt doesn't exist")
				return pi
			}

			pi.unit, ok = unitInterface.(string)
			if !ok {
				fmt.Printf("unit elemnt is not string")
				return pi
			}

			pricePerUnitInterface, exists := priceDimention["pricePerUnit"]
			if !exists {
				fmt.Printf("pricePerUnit elemnt doesn't exist")
				return pi
			}

			pricePerUnit, ok := pricePerUnitInterface.(map[string]interface{})
			if !ok {
				fmt.Printf("pricePerUnit elemnt is not map")
				return pi
			}

			for k, v := range pricePerUnit {
				pi.currency = k

				pi.valuePerUnit, ok = v.(string)
				if !ok {
					fmt.Printf("valuePerUnit elemnt is not map")
					return pi
				}
				return pi
			}

		}
	}

	return pi
}

func (u *uuNodesPlugin) getPrices() {
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

	// TODO bug in sdk?
	cfg.Region = "us-east-1"
	pricingService := pricing.New(cfg)

	productsInput := &pricing.GetProductsInput{
		Filters: []pricing.Filter{
			{
				Field: aws.String("ServiceCode"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("AmazonEC2"),
			},
			{
				Field: aws.String("productFamily"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Compute Instance"),
			},
			{
				Field: aws.String("operatingSystem"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Linux"),
			},
			{
				Field: aws.String("preInstalledSw"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("NA"),
			},
			//TODO: FIRST PRIORITY FIX, to filter by usagetype "EC2: Running Hours"
			//https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/selectdim.html
			//{
			//	Field: aws.String("tenancy"),
			//	Type:  pricing.FilterTypeTermMatch,
			//	Value: aws.String("Shared"),
			//},
			{
				Field: aws.String("location"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String(partitions[awsConfig.GetRegion()]), //TODO region to location??? bug, add PR to lib?
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(100), //TODO: add pagination
		ServiceCode:   aws.String("AmazonEC2"),
	}

	productsRequest := pricingService.GetProductsRequest(productsInput)

	productsPager := productsRequest.Paginate()
	for productsPager.Next() {
		page := productsPager.CurrentPage()

		if page != nil {
			for _, productItem := range page.PriceList {
				//b, _ := json.Marshal(productItem)
				var newPriceItem = getProduct(productItem)
				//TODO: some prices even for usagetype HostBoxUsage equal to zero. need to fix it later
				if newPriceItem.instanceType == "" || newPriceItem.memory == "" || newPriceItem.vcpu == "" || newPriceItem.valuePerUnit == "0.0000000000" {
					//b, _ := json.Marshal(productItem)
					//fmt.Printf("%s\n", b)
				}
				_, exists := u.computeInstancesPrices[newPriceItem.instanceType]
				if !exists {
					u.computeInstancesPrices[newPriceItem.instanceType] = make([]priceItem, 0, 0)
				}
				u.computeInstancesPrices[newPriceItem.instanceType] = append(u.computeInstancesPrices[newPriceItem.instanceType], newPriceItem)
			}
		}
	}

	if err = productsPager.Err(); err != nil {
		fmt.Printf("failed to describe products, %v", err)
	}

	fmt.Printf("found product prices: %v\n", len(u.computeInstancesPrices))
}

func (u *uuNodesPlugin) getInstanceEntries() (map[string]*InstanceEntry, error) {
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

	var instanceEntries = map[string]*InstanceEntry{}

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

		var entry = &InstanceEntry{
			AwsInstance: &AwsInstance{},
			KubeWorker: &KubeWorker{
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

func (u *uuNodesPlugin) getEC2Reservations() ([]ec2.RunInstancesOutput, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("unable to load AWS SDK config,  %v", err)
		return nil, errors.Wrap(err, "unable to load AWS SDK config")
	}

	var awsConfig = u.config.GetAwsConfig()

	cfg.Credentials = aws.NewStaticCredentialsProvider(
		awsConfig.GetAccessKeyId(),
		awsConfig.GetSecretAccessKey(),
		"",
	)
	cfg.Region = awsConfig.GetRegion()

	ec2Service := ec2.New(cfg)

	instancesRequest := ec2Service.DescribeInstancesRequest(nil)

	describeInstancesResponse, err := instancesRequest.Send()
	if err != nil {
		return nil, err
	}

	return describeInstancesResponse.Reservations, nil
}
