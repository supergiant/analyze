package requestslimitscheck

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/fields"

	"github.com/golang/protobuf/ptypes/any"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	corev1api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type resourceRequirementsPlugin struct {
	config *proto.PluginConfig
}

// TODO: this addition is till MVP, need to think and redesign while pluggability implementation
type checkResponse struct {
	description string
	Nodes       []nodeResourceRequirements
}

type nodeResourceRequirements struct {
	NodeName string
	Pods     []podResourceRequirements
}

type podResourceRequirements struct {
	PodName    string
	Containers []containerResourceRequirements
}

type containerResourceRequirements struct {
	ContainerName  string
	ContainerImage string
	Requests       struct {
		RAM string
		CPU string
	}
	Limits struct {
		RAM string
		CPU string
	}
}

const isNotSetStatus = "Is Not Set."

func NewPlugin() proto.PluginClient {
	return &resourceRequirementsPlugin{}
}

//TODO: wrap errors with meaningful messages
func (u *resourceRequirementsPlugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	var result = &proto.CheckResult{
		ExecutionStatus: "OK",
		Status:          proto.CheckStatus_UNKNOWN_CHECK_STATUS,
		Name:            "Resources (CPU/RAM) requests and limits Check",
		Description: &any.Any{
			TypeUrl: "io.supergiant.analyze.plugin.requestslimitscheck",
			Value:   nil,
		},
		Actions: []*proto.Action{
			&proto.Action{
				ActionId:    "1",
				Description: "Dismiss notification",
			},
			&proto.Action{
				ActionId:    "2",
				Description: "Set missing requests/limits",
			},
		},
	}

	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		result.ExecutionStatus = err.Error()

		return &proto.CheckResponse{Result: result}, nil
	}

	var descriptionValue = &checkResponse{
		description: "Resources (CPU/RAM) requests and limits where checked on nodes of k8s cluster.",
		Nodes:       make([]nodeResourceRequirements, 0, len(nodes.Items)),
	}

	for _, node := range nodes.Items {

		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name + ",status.phase!=" + string(corev1api.PodSucceeded) + ",status.phase!=" + string(corev1api.PodFailed))
		if err != nil {
			return nil, err
		}

		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: fieldSelector.String(),
		})
		if err != nil {
			if err != nil {
				result.ExecutionStatus = err.Error()

				return &proto.CheckResponse{Result: result}, nil
			}
		}
		var nodeDesc = nodeResourceRequirements{
			NodeName: node.Name,
			Pods:     make([]podResourceRequirements, 0, len(pods.Items)),
		}

		for _, pod := range pods.Items {
			var podDescription = podResourceRequirements{
				PodName:    pod.Name,
				Containers: make([]containerResourceRequirements, 0, len(pod.Spec.Containers)),
			}

			for _, container := range pod.Spec.Containers {
				resourceRequirementDescription, status := describeResourceRequirements(container)
				podDescription.Containers = append(podDescription.Containers, resourceRequirementDescription)
				setHigher(result, status)
			}
			nodeDesc.Pods = append(nodeDesc.Pods, podDescription)
		}

		descriptionValue.Nodes = append(descriptionValue.Nodes, nodeDesc)
	}

	bytes, err := json.Marshal(descriptionValue)
	if err != nil {
		return nil, err
	}

	result.Description.Value = bytes

	return &proto.CheckResponse{Result: result}, nil
}

func (u *resourceRequirementsPlugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *resourceRequirementsPlugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*proto.Empty, error) {
	return nil, nil
}

func (u *resourceRequirementsPlugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *resourceRequirementsPlugin) Info(ctx context.Context, in *proto.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
	return &proto.PluginInfo{
		Id:      "supergiant-resources-requests-and-limits-check-plugin",
		Version: "v0.0.1",
		Name:    "Resources (CPU/RAM) requests and limits",
		Description: "This plugin checks resources (CPU/RAM) requests and limits on pods. " +
			"It returns Green when limits and requests are set, Yellow when limits are not set, " +
			"and Red when requests are not set.",
	}, nil
}

func setHigher(ch *proto.CheckResult, status proto.CheckStatus) {
	if status > ch.Status {
		ch.Status = status
	}
}

func describeResourceRequirements(container v1.Container) (containerResourceRequirements, proto.CheckStatus) {
	var result = containerResourceRequirements{}
	result.ContainerName = container.Name
	result.ContainerImage = container.Image
	var resultStatus = proto.CheckStatus_GREEN
	var limitIsAbsent bool
	var requestIsAbsent bool

	if !container.Resources.Limits.Cpu().IsZero() {
		result.Limits.CPU = container.Resources.Limits.Cpu().String()
	} else {
		result.Limits.CPU = isNotSetStatus
		limitIsAbsent = true
	}

	if !container.Resources.Limits.Memory().IsZero() {
		result.Limits.RAM = container.Resources.Limits.Memory().String()
	} else {
		result.Limits.RAM = isNotSetStatus
		limitIsAbsent = true
	}

	if !container.Resources.Requests.Cpu().IsZero() {
		result.Requests.CPU = container.Resources.Requests.Cpu().String()
	} else {
		result.Requests.CPU = isNotSetStatus
		requestIsAbsent = true
	}

	if !container.Resources.Requests.Memory().IsZero() {
		result.Requests.RAM = container.Resources.Requests.Memory().String()
	} else {
		result.Requests.RAM = isNotSetStatus
		requestIsAbsent = true
	}

	if limitIsAbsent {
		resultStatus = proto.CheckStatus_YELLOW
	}

	if requestIsAbsent {
		resultStatus = proto.CheckStatus_RED
	}

	return result, resultStatus
}
