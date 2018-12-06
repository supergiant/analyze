package sunsetting

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/supergiant/analyze/builtin/plugins/sunsetting/cloudprovider"
	"github.com/supergiant/analyze/builtin/plugins/sunsetting/cloudprovider/aws"
	"github.com/supergiant/analyze/builtin/plugins/sunsetting/kube"
	"github.com/supergiant/analyze/builtin/plugins/sunsetting/nodeagent"
	"github.com/supergiant/analyze/pkg/plugin/proto"
)

type plugin struct {
	config                 *proto.PluginConfig
	nodeAgentClient        *nodeagent.Client
	awsClient              *aws.Client
	kubeClient             *kube.Client
	computeInstancesPrices map[string][]cloudprovider.ProductPrice
	logger                 logrus.FieldLogger
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
			Name:        "Dismiss notification",
			Description: "Dismiss notification, just prevents notification from being shown",
		},
		&proto.Action{
			ActionId:    "2",
			Name:        "Sunset nodes",
			Description: "Sunset nodes, makes request to capacity service to remove underutilized nodes.",
		},
	},
}

func NewPlugin() proto.PluginClient {
	return &plugin{
		logger: logrus.New(),
	}
}

func (u *plugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
	var nodeResourceRequirements, err = u.kubeClient.GetNodeResourceRequirements()
	if err != nil {
		fmt.Printf("unable to get nodeResourceRequirements, %v", err)
		return nil, errors.Wrap(err, "unable to get nodeResourceRequirements")
	}

	var computeInstances = make(map[string]cloudprovider.ComputeInstance)
	for instanceID, resourceRequirements := range nodeResourceRequirements {
		var nodeAgentInstance = nodeagent.Instance{URI: resourceRequirements.IPAddress()}
		fetchedInstanceID, err := u.nodeAgentClient.Get(nodeAgentInstance, "/aws/meta-data/instance-id")
		if err != nil {
			u.logger.Errorf("cant fetch ID for node %s", instanceID)
			continue
		}
		if fetchedInstanceID != instanceID {
			u.logger.Errorf(
				"fetched ec2 instanceID: %s not equal to instanceID from providerID %s",
				fetchedInstanceID,
				instanceID,
			)
			continue
		}

		instanceType, err := u.nodeAgentClient.Get(nodeAgentInstance, "/aws/meta-data/instance-type")
		computeInstances[instanceID] = cloudprovider.ComputeInstance{
			InstanceID:   instanceID,
			InstanceType: instanceType,
		}
	}

	//computeInstances, err := u.awsClient.GetComputeInstances()
	//if err != nil {
	//	fmt.Printf("failed to describe ec2 instances, %v", err)
	//	return nil, errors.Wrap(err, "failed to describe ec2 instances")
	//}

	var unsortedEntries []*InstanceEntry
	var result []InstanceEntry

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
				if strings.Contains(priceItem.UsageType, "BoxUsage") && priceItem.ValuePerUnit != "0.0000000000" {
					instanceTypePrice = priceItem
				}
			}
			if instanceTypePrice.InstanceType == "" && len(instanceTypePrices) > 0 {
				instanceTypePrice = instanceTypePrices[0]
			}
		}

		result = append(result, InstanceEntry{
			CloudProvider: computeInstance,
			Price:         instanceTypePrice,
			WorkerNode:    *kubeNode,
		})
		unsortedEntries = append(unsortedEntries, &InstanceEntry{
			CloudProvider: computeInstance,
			Price:         instanceTypePrice,
			WorkerNode:    *kubeNode,
		})
	}

	//TODO: double check logic, is it really needed?
	var instancesToSunset = CheckEachPodOneByOne(unsortedEntries)
	if len(instancesToSunset) == 0 {
		instancesToSunset = CheckAllPodsAtATime(unsortedEntries)
	}

	// mark nodes selected node with IsRecommendedToSunset == true
	for i, _ := range result {
		for _, entryToSunset := range instancesToSunset {
			if entryToSunset.CloudProvider.InstanceID == result[i].CloudProvider.InstanceID {
				result[i].WorkerNode.IsRecommendedToSunset = true
			}
		}
	}

	b, _ := json.Marshal(result)

	checkResult.Description = &any.Any{
		TypeUrl: "io.supergiant.analyze.plugin.sunsetting",
		Value:   b,
	}

	if len(instancesToSunset) == 0 {
		checkResult.Status = proto.CheckStatus_GREEN
	} else {
		checkResult.Status = proto.CheckStatus_YELLOW
	}

	return &proto.CheckResponse{Result: checkResult}, nil
}

func (u *plugin) Configure(ctx context.Context, pluginConfig *proto.PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	//TODO: add here config validation in future
	u.config = pluginConfig

	nodeAgentClient, err := nodeagent.NewClient(logrus.New())
	if err != nil {
		return nil, err
	}
	u.nodeAgentClient = nodeAgentClient

	var awsClientConfig = pluginConfig.GetAwsConfig()
	awsClient, err := aws.NewClient(awsClientConfig)
	if err != nil {
		return nil, err
	}
	u.awsClient = awsClient

	//TODO: may be we need just log warning?
	u.computeInstancesPrices, err = u.awsClient.GetPrices()
	if err != nil {
		return nil, err
	}

	u.kubeClient, err = kube.NewKubeClient()
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
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
