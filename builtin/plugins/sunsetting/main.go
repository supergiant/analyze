package sunsetting

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/builtin/plugins/sunsetting/kube"
	"github.com/supergiant/robot/builtin/plugins/sunsetting/prices"
	"github.com/supergiant/robot/pkg/plugin/proto"
)

type plugin struct {
	config       *proto.PluginConfig
	ec2Service   *ec2.EC2
	сoreV1Client *corev1.CoreV1Client

	computeInstancesPrices map[string][]prices.Item
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

	ec2Reservations, err := getEC2Reservations(u.ec2Service)
	if err != nil {
		fmt.Printf("failed to describe ec2 instances, %v", err)
		return nil, errors.Wrap(err, "failed to describe ec2 instances")
	}

	var unsortedEntries = []*InstanceEntry{}

	// create InstanceEntries by combining nodeResourceRequirements with ec2 instance type
	for _, instancesReservation := range ec2Reservations {
		for _, i := range instancesReservation.Instances {
			if i.InstanceId == nil {
				continue
			}
			var kubeNode, exists = nodeResourceRequirements[*i.InstanceId]
			if !exists {
				continue
			}

			var instanceType, _ = i.InstanceType.MarshalValue()

			unsortedEntries = append(unsortedEntries, &InstanceEntry{
				InstanceType:             instanceType,
				NodeResourceRequirements: kubeNode,
			})
		}
	}

	var sortedByWastedRam = NewSortedEntriesByWastedRAM(unsortedEntries)
	//var sortedByRequestedRam = models.NewSortedEntriesByRequestedRAM(unsortedEntries)

	var instancesToSunset = CheckAllPodsAtATime(sortedByWastedRam)

	b, _ := json.Marshal(instancesToSunset)

	checkResult.Description = &any.Any{
		TypeUrl: "test",
		Value:   b,
	}
	checkResult.Status = proto.CheckStatus_GREEN
	return &proto.CheckResponse{Result: checkResult}, nil
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

func getEC2Reservations(ec2Service *ec2.EC2) ([]ec2.RunInstancesOutput, error) {

	instancesRequest := ec2Service.DescribeInstancesRequest(nil)

	describeInstancesResponse, err := instancesRequest.Send()
	if err != nil {
		return nil, err
	}

	return describeInstancesResponse.Reservations, nil
}
