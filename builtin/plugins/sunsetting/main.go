package sunsetting

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/cloudprovider"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/cloudprovider/aws"
	"github.com/supergiant/robot/builtin/plugins/sunsetting/kube"
	"github.com/supergiant/robot/pkg/plugin/proto"
)

type plugin struct {
	config                 *proto.PluginConfig
	сoreV1Client           *corev1.CoreV1Client
	awsClient              *aws.Client
	computeInstancesPrices map[string][]cloudprovider.ProductPrice
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
	var nodeResourceRequirements, err = kube.GetNodeResourceRequirements(u.сoreV1Client)
	if err != nil {
		fmt.Printf("unable to get nodeResourceRequirements, %v", err)
		return nil, errors.Wrap(err, "unable to get nodeResourceRequirements")
	}

	computeInstances, err := u.awsClient.GetComputeInstances()
	if err != nil {
		fmt.Printf("failed to describe ec2 instances, %v", err)
		return nil, errors.Wrap(err, "failed to describe ec2 instances")
	}

	var unsortedEntries = []*InstanceEntry{}
	var unsorted = []InstanceEntry{}

	// create InstanceEntries by combining nodeResourceRequirements with ec2 instance type and price
	for InstanceID, computeInstance := range computeInstances {
		var kubeNode, exists = nodeResourceRequirements[InstanceID]
		if !exists {
			continue
		}

		// TODO: fix me when prices collecting will be clear
		// TODO: We need to match it with instance tenancy?
		var instanceTypePrice cloudprovider.ProductPrice
		var instanceTypePrices, exist = u.computeInstancesPrices[computeInstance.InstanceType]
		if exist {
			for _, priceItem := range instanceTypePrices {
				if strings.Contains(priceItem.UsageType, "BoxUsage") {
					instanceTypePrice = priceItem
				}
			}
			if instanceTypePrice.InstanceType == "" && len(instanceTypePrices) > 0 {
				instanceTypePrice = instanceTypePrices[0]
			}
		}

		unsorted = append(unsorted, InstanceEntry{
			CloudProvider:            computeInstance,
			Price:                    instanceTypePrice,
			NodeResourceRequirements: *kubeNode,
		})
		unsortedEntries = append(unsortedEntries, &InstanceEntry{
			CloudProvider:            computeInstance,
			Price:                    instanceTypePrice,
			NodeResourceRequirements: *kubeNode,
		})
	}

	var sortedByWastedRam = NewSortedEntriesByWastedRAM(unsortedEntries)
	var sortedByRequestedRam = NewSortedEntriesByRequestedRAM(unsortedEntries)

	var instancesToSunset = CheckAllPodsAtATime(sortedByWastedRam)

	var instancesToSunsetOptionTwo = CheckEachPodOneByOne(sortedByWastedRam, sortedByRequestedRam)

	var result = map[string]interface{}{
		"allInstances":               unsorted,
		"instancesToSunset":          instancesToSunset,
		"instancesToSunsetOptionTwo": instancesToSunsetOptionTwo,
	}
	b, _ := json.Marshal(result)

	checkResult.Description = &any.Any{
		TypeUrl: "io.supergiant.analyze.plugin.sunsetting",
		Value:   b,
	}
	checkResult.Status = proto.CheckStatus_GREEN
	return &proto.CheckResponse{Result: checkResult}, nil
}

func (u *plugin) Configure(ctx context.Context, pluginConfig *proto.PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	u.config = pluginConfig
	//TODO: add here config validation in future
	var awsClientConfig = pluginConfig.GetAwsConfig()
	var awsClient, err = aws.NewClient(awsClientConfig)
	if err != nil {
		return nil, err
	}
	u.awsClient = awsClient

	//TODO: may be we need just log warning?
	u.computeInstancesPrices, err = u.awsClient.GetPrices()
	if err != nil {
		return nil, err
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the client
	u.сoreV1Client, err = corev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *plugin) Info(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
	return &proto.PluginInfo{
		Id:          "supergiant-underutilized-nodes-plugin",
		Version:     "v0.0.1",
		Name:        "Underutilized nodes sunsetting plugin",
		Description: "This plugin checks nodes using high intelligent Kelly's approach to find underutilized nodes, than calculates how it is possible to fix that",
	}, nil
}

func (u *plugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *plugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}
