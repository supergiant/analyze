// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/plugin.proto

package proto

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	duration "github.com/golang/protobuf/ptypes/duration"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CheckStatus int32

const (
	CheckStatus_UNKNOWN_CHECK_STATUS CheckStatus = 0
	CheckStatus_GREEN                CheckStatus = 1
	CheckStatus_YELLOW               CheckStatus = 2
	CheckStatus_RED                  CheckStatus = 3
)

var CheckStatus_name = map[int32]string{
	0: "UNKNOWN_CHECK_STATUS",
	1: "GREEN",
	2: "YELLOW",
	3: "RED",
}

var CheckStatus_value = map[string]int32{
	"UNKNOWN_CHECK_STATUS": 0,
	"GREEN":                1,
	"YELLOW":               2,
	"RED":                  3,
}

func (x CheckStatus) String() string {
	return proto.EnumName(CheckStatus_name, int32(x))
}

func (CheckStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{0}
}

type CloudProviderType int32

const (
	CloudProviderType_UNKNOWN_CLOUD_PROVIDER_TYPE CloudProviderType = 0
	CloudProviderType_AWS                         CloudProviderType = 1
	CloudProviderType_DO                          CloudProviderType = 2
)

var CloudProviderType_name = map[int32]string{
	0: "UNKNOWN_CLOUD_PROVIDER_TYPE",
	1: "AWS",
	2: "DO",
}

var CloudProviderType_value = map[string]int32{
	"UNKNOWN_CLOUD_PROVIDER_TYPE": 0,
	"AWS":                         1,
	"DO":                          2,
}

func (x CloudProviderType) String() string {
	return proto.EnumName(CloudProviderType_name, int32(x))
}

func (CloudProviderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{1}
}

type PluginInfo struct {
	//unique ID of plugin
	//basically it is slugged URI of plugin repository name e. g. supergiant-underutilized-nodes
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	//plugin version in semver format
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	//short name of plugin
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// detailed plugin description
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// plugin default config
	DefaultConfig        *PluginConfig `protobuf:"bytes,5,opt,name=default_config,json=defaultConfig,proto3" json:"default_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PluginInfo) Reset()         { *m = PluginInfo{} }
func (m *PluginInfo) String() string { return proto.CompactTextString(m) }
func (*PluginInfo) ProtoMessage()    {}
func (*PluginInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{0}
}

func (m *PluginInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PluginInfo.Unmarshal(m, b)
}
func (m *PluginInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PluginInfo.Marshal(b, m, deterministic)
}
func (m *PluginInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginInfo.Merge(m, src)
}
func (m *PluginInfo) XXX_Size() int {
	return xxx_messageInfo_PluginInfo.Size(m)
}
func (m *PluginInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PluginInfo proto.InternalMessageInfo

func (m *PluginInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PluginInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *PluginInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PluginInfo) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *PluginInfo) GetDefaultConfig() *PluginConfig {
	if m != nil {
		return m.DefaultConfig
	}
	return nil
}

// plugin instance is running inside k8s cluster, consequently:
// k8s api server doesn't need to be configured and is accessible through k8s service discovery
// AWS credentials need to be configured
type PluginConfig struct {
	ExecutionInterval *duration.Duration `protobuf:"bytes,1,opt,name=execution_interval,json=executionInterval,proto3" json:"execution_interval,omitempty"`
	// must be valid JSON
	PluginSpecificConfig []byte            `protobuf:"bytes,2,opt,name=plugin_specific_config,json=pluginSpecificConfig,proto3" json:"plugin_specific_config,omitempty"`
	ProviderType         CloudProviderType `protobuf:"varint,3,opt,name=provider_type,json=providerType,proto3,enum=proto.CloudProviderType" json:"provider_type,omitempty"`
	// Types that are valid to be assigned to CloudProviderConfig:
	//	*PluginConfig_DoConfig
	//	*PluginConfig_AwsConfig
	CloudProviderConfig  isPluginConfig_CloudProviderConfig `protobuf_oneof:"cloud_provider_config"`
	MetricsServerUri     string                             `protobuf:"bytes,6,opt,name=metrics_server_uri,json=metricsServerUri,proto3" json:"metrics_server_uri,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *PluginConfig) Reset()         { *m = PluginConfig{} }
func (m *PluginConfig) String() string { return proto.CompactTextString(m) }
func (*PluginConfig) ProtoMessage()    {}
func (*PluginConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{1}
}

func (m *PluginConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PluginConfig.Unmarshal(m, b)
}
func (m *PluginConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PluginConfig.Marshal(b, m, deterministic)
}
func (m *PluginConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginConfig.Merge(m, src)
}
func (m *PluginConfig) XXX_Size() int {
	return xxx_messageInfo_PluginConfig.Size(m)
}
func (m *PluginConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginConfig.DiscardUnknown(m)
}

var xxx_messageInfo_PluginConfig proto.InternalMessageInfo

func (m *PluginConfig) GetExecutionInterval() *duration.Duration {
	if m != nil {
		return m.ExecutionInterval
	}
	return nil
}

func (m *PluginConfig) GetPluginSpecificConfig() []byte {
	if m != nil {
		return m.PluginSpecificConfig
	}
	return nil
}

func (m *PluginConfig) GetProviderType() CloudProviderType {
	if m != nil {
		return m.ProviderType
	}
	return CloudProviderType_UNKNOWN_CLOUD_PROVIDER_TYPE
}

type isPluginConfig_CloudProviderConfig interface {
	isPluginConfig_CloudProviderConfig()
}

type PluginConfig_DoConfig struct {
	DoConfig *DoConfig `protobuf:"bytes,4,opt,name=do_config,json=doConfig,proto3,oneof"`
}

type PluginConfig_AwsConfig struct {
	AwsConfig *AwsConfig `protobuf:"bytes,5,opt,name=aws_config,json=awsConfig,proto3,oneof"`
}

func (*PluginConfig_DoConfig) isPluginConfig_CloudProviderConfig() {}

func (*PluginConfig_AwsConfig) isPluginConfig_CloudProviderConfig() {}

func (m *PluginConfig) GetCloudProviderConfig() isPluginConfig_CloudProviderConfig {
	if m != nil {
		return m.CloudProviderConfig
	}
	return nil
}

func (m *PluginConfig) GetDoConfig() *DoConfig {
	if x, ok := m.GetCloudProviderConfig().(*PluginConfig_DoConfig); ok {
		return x.DoConfig
	}
	return nil
}

func (m *PluginConfig) GetAwsConfig() *AwsConfig {
	if x, ok := m.GetCloudProviderConfig().(*PluginConfig_AwsConfig); ok {
		return x.AwsConfig
	}
	return nil
}

func (m *PluginConfig) GetMetricsServerUri() string {
	if m != nil {
		return m.MetricsServerUri
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*PluginConfig) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*PluginConfig_DoConfig)(nil),
		(*PluginConfig_AwsConfig)(nil),
	}
}

type Stop struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stop) Reset()         { *m = Stop{} }
func (m *Stop) String() string { return proto.CompactTextString(m) }
func (*Stop) ProtoMessage()    {}
func (*Stop) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{2}
}

func (m *Stop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stop.Unmarshal(m, b)
}
func (m *Stop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stop.Marshal(b, m, deterministic)
}
func (m *Stop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stop.Merge(m, src)
}
func (m *Stop) XXX_Size() int {
	return xxx_messageInfo_Stop.Size(m)
}
func (m *Stop) XXX_DiscardUnknown() {
	xxx_messageInfo_Stop.DiscardUnknown(m)
}

var xxx_messageInfo_Stop proto.InternalMessageInfo

type Stop_Request struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stop_Request) Reset()         { *m = Stop_Request{} }
func (m *Stop_Request) String() string { return proto.CompactTextString(m) }
func (*Stop_Request) ProtoMessage()    {}
func (*Stop_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{2, 0}
}

func (m *Stop_Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stop_Request.Unmarshal(m, b)
}
func (m *Stop_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stop_Request.Marshal(b, m, deterministic)
}
func (m *Stop_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stop_Request.Merge(m, src)
}
func (m *Stop_Request) XXX_Size() int {
	return xxx_messageInfo_Stop_Request.Size(m)
}
func (m *Stop_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Stop_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Stop_Request proto.InternalMessageInfo

type Stop_Response struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stop_Response) Reset()         { *m = Stop_Response{} }
func (m *Stop_Response) String() string { return proto.CompactTextString(m) }
func (*Stop_Response) ProtoMessage()    {}
func (*Stop_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{2, 1}
}

func (m *Stop_Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stop_Response.Unmarshal(m, b)
}
func (m *Stop_Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stop_Response.Marshal(b, m, deterministic)
}
func (m *Stop_Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stop_Response.Merge(m, src)
}
func (m *Stop_Response) XXX_Size() int {
	return xxx_messageInfo_Stop_Response.Size(m)
}
func (m *Stop_Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Stop_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Stop_Response proto.InternalMessageInfo

func (m *Stop_Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type CheckRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckRequest) Reset()         { *m = CheckRequest{} }
func (m *CheckRequest) String() string { return proto.CompactTextString(m) }
func (*CheckRequest) ProtoMessage()    {}
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{3}
}

func (m *CheckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckRequest.Unmarshal(m, b)
}
func (m *CheckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckRequest.Marshal(b, m, deterministic)
}
func (m *CheckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckRequest.Merge(m, src)
}
func (m *CheckRequest) XXX_Size() int {
	return xxx_messageInfo_CheckRequest.Size(m)
}
func (m *CheckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckRequest proto.InternalMessageInfo

type CheckResponse struct {
	Result               *CheckResult `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Error                string       `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CheckResponse) Reset()         { *m = CheckResponse{} }
func (m *CheckResponse) String() string { return proto.CompactTextString(m) }
func (*CheckResponse) ProtoMessage()    {}
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{4}
}

func (m *CheckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckResponse.Unmarshal(m, b)
}
func (m *CheckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckResponse.Marshal(b, m, deterministic)
}
func (m *CheckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckResponse.Merge(m, src)
}
func (m *CheckResponse) XXX_Size() int {
	return xxx_messageInfo_CheckResponse.Size(m)
}
func (m *CheckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CheckResponse proto.InternalMessageInfo

func (m *CheckResponse) GetResult() *CheckResult {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *CheckResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type CheckResult struct {
	ExecutionStatus string      `protobuf:"bytes,1,opt,name=execution_status,json=executionStatus,proto3" json:"execution_status,omitempty"`
	Status          CheckStatus `protobuf:"varint,2,opt,name=status,proto3,enum=proto.CheckStatus" json:"status,omitempty"`
	Name            string      `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// TODO: temporary solution need to be redesigned.
	Description          *any.Any `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckResult) Reset()         { *m = CheckResult{} }
func (m *CheckResult) String() string { return proto.CompactTextString(m) }
func (*CheckResult) ProtoMessage()    {}
func (*CheckResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{5}
}

func (m *CheckResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckResult.Unmarshal(m, b)
}
func (m *CheckResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckResult.Marshal(b, m, deterministic)
}
func (m *CheckResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckResult.Merge(m, src)
}
func (m *CheckResult) XXX_Size() int {
	return xxx_messageInfo_CheckResult.Size(m)
}
func (m *CheckResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckResult.DiscardUnknown(m)
}

var xxx_messageInfo_CheckResult proto.InternalMessageInfo

func (m *CheckResult) GetExecutionStatus() string {
	if m != nil {
		return m.ExecutionStatus
	}
	return ""
}

func (m *CheckResult) GetStatus() CheckStatus {
	if m != nil {
		return m.Status
	}
	return CheckStatus_UNKNOWN_CHECK_STATUS
}

func (m *CheckResult) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CheckResult) GetDescription() *any.Any {
	if m != nil {
		return m.Description
	}
	return nil
}

type AwsConfig struct {
	AccessKeyId          string   `protobuf:"bytes,1,opt,name=access_key_id,json=accessKeyId,proto3" json:"access_key_id,omitempty"`
	SecretAccessKey      string   `protobuf:"bytes,2,opt,name=secret_access_key,json=secretAccessKey,proto3" json:"secret_access_key,omitempty"`
	Region               string   `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AwsConfig) Reset()         { *m = AwsConfig{} }
func (m *AwsConfig) String() string { return proto.CompactTextString(m) }
func (*AwsConfig) ProtoMessage()    {}
func (*AwsConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{6}
}

func (m *AwsConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AwsConfig.Unmarshal(m, b)
}
func (m *AwsConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AwsConfig.Marshal(b, m, deterministic)
}
func (m *AwsConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AwsConfig.Merge(m, src)
}
func (m *AwsConfig) XXX_Size() int {
	return xxx_messageInfo_AwsConfig.Size(m)
}
func (m *AwsConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AwsConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AwsConfig proto.InternalMessageInfo

func (m *AwsConfig) GetAccessKeyId() string {
	if m != nil {
		return m.AccessKeyId
	}
	return ""
}

func (m *AwsConfig) GetSecretAccessKey() string {
	if m != nil {
		return m.SecretAccessKey
	}
	return ""
}

func (m *AwsConfig) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

// not implemented yet
type DoConfig struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DoConfig) Reset()         { *m = DoConfig{} }
func (m *DoConfig) String() string { return proto.CompactTextString(m) }
func (*DoConfig) ProtoMessage()    {}
func (*DoConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_4368bc2f41172306, []int{7}
}

func (m *DoConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DoConfig.Unmarshal(m, b)
}
func (m *DoConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DoConfig.Marshal(b, m, deterministic)
}
func (m *DoConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DoConfig.Merge(m, src)
}
func (m *DoConfig) XXX_Size() int {
	return xxx_messageInfo_DoConfig.Size(m)
}
func (m *DoConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DoConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DoConfig proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("proto.CheckStatus", CheckStatus_name, CheckStatus_value)
	proto.RegisterEnum("proto.CloudProviderType", CloudProviderType_name, CloudProviderType_value)
	proto.RegisterType((*PluginInfo)(nil), "proto.PluginInfo")
	proto.RegisterType((*PluginConfig)(nil), "proto.PluginConfig")
	proto.RegisterType((*Stop)(nil), "proto.Stop")
	proto.RegisterType((*Stop_Request)(nil), "proto.Stop.Request")
	proto.RegisterType((*Stop_Response)(nil), "proto.Stop.Response")
	proto.RegisterType((*CheckRequest)(nil), "proto.CheckRequest")
	proto.RegisterType((*CheckResponse)(nil), "proto.CheckResponse")
	proto.RegisterType((*CheckResult)(nil), "proto.CheckResult")
	proto.RegisterType((*AwsConfig)(nil), "proto.AwsConfig")
	proto.RegisterType((*DoConfig)(nil), "proto.DoConfig")
}

func init() { proto.RegisterFile("proto/plugin.proto", fileDescriptor_4368bc2f41172306) }

var fileDescriptor_4368bc2f41172306 = []byte{
	// 749 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0xdd, 0x6e, 0xda, 0x4c,
	0x10, 0x8d, 0x1d, 0x20, 0x61, 0xf8, 0x89, 0xd9, 0xf0, 0xe5, 0x73, 0x88, 0xd4, 0x22, 0x5f, 0xa5,
	0xa8, 0x22, 0x0a, 0xa9, 0xaa, 0xaa, 0x52, 0x2f, 0x28, 0x58, 0x09, 0x4a, 0x04, 0xd4, 0x40, 0xa3,
	0x5c, 0x59, 0x8e, 0xbd, 0x50, 0x2b, 0xc4, 0x76, 0x77, 0x6d, 0x52, 0xde, 0xa8, 0x97, 0x7d, 0xa0,
	0xbe, 0x47, 0x6f, 0x2b, 0xef, 0xae, 0x1d, 0x7e, 0xd2, 0x2b, 0x98, 0x99, 0x33, 0x3e, 0x67, 0x66,
	0xcf, 0x2e, 0xa0, 0x80, 0xf8, 0xa1, 0x7f, 0x16, 0xcc, 0xa3, 0x99, 0xeb, 0x35, 0x59, 0x80, 0xb2,
	0xec, 0xa7, 0x76, 0x3c, 0xf3, 0xfd, 0xd9, 0x1c, 0x9f, 0xb1, 0xe8, 0x3e, 0x9a, 0x9e, 0x59, 0xde,
	0x92, 0x23, 0x6a, 0x27, 0x9b, 0x25, 0xfc, 0x18, 0x84, 0x49, 0xf1, 0xd5, 0x66, 0xd1, 0x89, 0x88,
	0x15, 0xba, 0xbe, 0xf8, 0xbc, 0xf6, 0x53, 0x02, 0x18, 0x32, 0xbe, 0x9e, 0x37, 0xf5, 0x51, 0x19,
	0x64, 0xd7, 0x51, 0xa5, 0xba, 0x74, 0x9a, 0x37, 0x64, 0xd7, 0x41, 0x2a, 0xec, 0x2d, 0x30, 0xa1,
	0xae, 0xef, 0xa9, 0x32, 0x4b, 0x26, 0x21, 0x42, 0x90, 0xf1, 0xac, 0x47, 0xac, 0xee, 0xb2, 0x34,
	0xfb, 0x8f, 0xea, 0x50, 0x70, 0x30, 0xb5, 0x89, 0x1b, 0xc4, 0x0c, 0x6a, 0x86, 0x95, 0x56, 0x53,
	0xe8, 0x23, 0x94, 0x1d, 0x3c, 0xb5, 0xa2, 0x79, 0x68, 0xda, 0xbe, 0x37, 0x75, 0x67, 0x6a, 0xb6,
	0x2e, 0x9d, 0x16, 0x5a, 0x87, 0x5c, 0x4e, 0x93, 0x4b, 0xe9, 0xb0, 0x92, 0x51, 0x12, 0x50, 0x1e,
	0x6a, 0x7f, 0x64, 0x28, 0xae, 0xd6, 0xd1, 0x15, 0x20, 0xfc, 0x03, 0xdb, 0x51, 0xfc, 0x65, 0xd3,
	0xf5, 0x42, 0x4c, 0x16, 0xd6, 0x9c, 0x89, 0x2f, 0xb4, 0x8e, 0x9b, 0x7c, 0xf0, 0x66, 0x32, 0x78,
	0xb3, 0x2b, 0x06, 0x37, 0x2a, 0x69, 0x53, 0x4f, 0xf4, 0xa0, 0x77, 0x70, 0xc4, 0x97, 0x6e, 0xd2,
	0x00, 0xdb, 0xee, 0xd4, 0xb5, 0x13, 0x79, 0xf1, 0xd4, 0x45, 0xa3, 0xca, 0xab, 0x23, 0x51, 0x14,
	0xfc, 0x9f, 0xa0, 0x14, 0x10, 0x7f, 0xe1, 0x3a, 0x98, 0x98, 0xe1, 0x32, 0xe0, 0xbb, 0x28, 0xb7,
	0x54, 0x31, 0x4b, 0x67, 0xee, 0x47, 0xce, 0x50, 0x00, 0xc6, 0xcb, 0x00, 0x1b, 0xc5, 0x60, 0x25,
	0x42, 0x4d, 0xc8, 0x3b, 0x7e, 0xc2, 0x93, 0x61, 0xaa, 0x0f, 0x44, 0x6b, 0xd7, 0xe7, 0x14, 0x57,
	0x3b, 0xc6, 0xbe, 0x23, 0xfe, 0xa3, 0x73, 0x00, 0xeb, 0x89, 0xae, 0xef, 0x4d, 0x11, 0x0d, 0xed,
	0x27, 0x9a, 0x76, 0xe4, 0xad, 0x24, 0x40, 0x6f, 0x01, 0x3d, 0xe2, 0x90, 0xb8, 0x36, 0x35, 0x29,
	0x26, 0x0b, 0x4c, 0xcc, 0x88, 0xb8, 0x6a, 0x8e, 0x9d, 0x8b, 0x22, 0x2a, 0x23, 0x56, 0x98, 0x10,
	0xf7, 0xf3, 0xff, 0xf0, 0x9f, 0x1d, 0x6b, 0x36, 0xd3, 0xa9, 0x38, 0x97, 0x76, 0x01, 0x99, 0x51,
	0xe8, 0x07, 0xb5, 0x3c, 0xec, 0x19, 0xf8, 0x7b, 0x84, 0x69, 0x58, 0xab, 0xc3, 0xbe, 0x81, 0x69,
	0xe0, 0x7b, 0x14, 0xa3, 0x2a, 0x64, 0x31, 0x21, 0x3e, 0x11, 0xbe, 0xe1, 0x81, 0x56, 0x86, 0x62,
	0xe7, 0x1b, 0xb6, 0x1f, 0x44, 0x87, 0xf6, 0x05, 0x4a, 0x22, 0x16, 0x6d, 0x0d, 0xc8, 0x11, 0x4c,
	0xa3, 0x79, 0x28, 0x8e, 0x0c, 0x25, 0x7b, 0x13, 0xa8, 0x68, 0x1e, 0x1a, 0x02, 0xf1, 0x4c, 0x21,
	0xaf, 0x52, 0xfc, 0x92, 0xa0, 0xb0, 0x82, 0x46, 0x6f, 0x40, 0x79, 0x36, 0x04, 0x0d, 0xad, 0x30,
	0xa2, 0x42, 0xd3, 0x41, 0x9a, 0x1f, 0xb1, 0x74, 0x4c, 0x2e, 0x00, 0x32, 0x3b, 0xb4, 0x35, 0x72,
	0x8e, 0x31, 0x04, 0xe2, 0x45, 0xab, 0xbf, 0xdf, 0xb6, 0x7a, 0xa1, 0x55, 0xdd, 0x32, 0x5d, 0xdb,
	0x5b, 0xae, 0x5d, 0x00, 0x8d, 0x42, 0x3e, 0x3d, 0x2b, 0xa4, 0x41, 0xc9, 0xb2, 0x6d, 0x4c, 0xa9,
	0xf9, 0x80, 0x97, 0x66, 0x7a, 0xf1, 0x0a, 0x3c, 0x79, 0x8d, 0x97, 0x3d, 0x07, 0x35, 0xa0, 0x42,
	0xb1, 0x4d, 0x70, 0x68, 0x3e, 0x43, 0xc5, 0x16, 0x0e, 0x78, 0xa1, 0x9d, 0xa0, 0xd1, 0x51, 0xbc,
	0xd1, 0x59, 0xac, 0x87, 0x4b, 0x15, 0x91, 0x06, 0xb0, 0x9f, 0x38, 0xaa, 0x71, 0x29, 0x56, 0x26,
	0xf6, 0xa0, 0x42, 0x75, 0xd2, 0xbf, 0xee, 0x0f, 0x6e, 0xfb, 0x66, 0xe7, 0x4a, 0xef, 0x5c, 0x9b,
	0xa3, 0x71, 0x7b, 0x3c, 0x19, 0x29, 0x3b, 0x28, 0x0f, 0xd9, 0x4b, 0x43, 0xd7, 0xfb, 0x8a, 0x84,
	0x00, 0x72, 0x77, 0xfa, 0xcd, 0xcd, 0xe0, 0x56, 0x91, 0xd1, 0x1e, 0xec, 0x1a, 0x7a, 0x57, 0xd9,
	0x6d, 0xe8, 0x50, 0xd9, 0x72, 0x38, 0x7a, 0x0d, 0x27, 0xe9, 0xe7, 0x6e, 0x06, 0x93, 0xae, 0x39,
	0x34, 0x06, 0x5f, 0x7b, 0x5d, 0xdd, 0x30, 0xc7, 0x77, 0x43, 0x5d, 0xd9, 0x89, 0xdb, 0xdb, 0xb7,
	0x23, 0x45, 0x42, 0x39, 0x90, 0xbb, 0x03, 0x45, 0x6e, 0xfd, 0x96, 0x20, 0xc7, 0x6f, 0x35, 0x3a,
	0x87, 0x0c, 0x7b, 0x84, 0x8e, 0xb6, 0xd6, 0xa8, 0xc7, 0x2f, 0x5a, 0xad, 0xb2, 0xf6, 0x48, 0x30,
	0x68, 0x0b, 0xb2, 0x6c, 0x1a, 0x74, 0xb8, 0x6e, 0x1e, 0x6e, 0xd2, 0xea, 0x86, 0xa3, 0xb8, 0xef,
	0x3e, 0x40, 0x9e, 0xef, 0x22, 0x22, 0x18, 0xbd, 0xf4, 0xf0, 0xd4, 0xfe, 0x21, 0x20, 0x16, 0x18,
	0xdf, 0x83, 0xb4, 0x29, 0x0e, 0x9a, 0x9b, 0x64, 0x22, 0xc9, 0xc9, 0xee, 0x73, 0x2c, 0x79, 0xf1,
	0x37, 0x00, 0x00, 0xff, 0xff, 0x59, 0x3e, 0xc0, 0x10, 0xdb, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PluginClient interface {
	Info(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PluginInfo, error)
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	Configure(ctx context.Context, in *PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error)
	//Graceful Shutdown Plugin
	Stop(ctx context.Context, in *Stop_Request, opts ...grpc.CallOption) (*Stop_Response, error)
}

type pluginClient struct {
	cc *grpc.ClientConn
}

func NewPluginClient(cc *grpc.ClientConn) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) Info(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PluginInfo, error) {
	out := new(PluginInfo)
	err := c.cc.Invoke(ctx, "/proto.Plugin/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/proto.Plugin/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Configure(ctx context.Context, in *PluginConfig, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.Plugin/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Stop(ctx context.Context, in *Stop_Request, opts ...grpc.CallOption) (*Stop_Response, error) {
	out := new(Stop_Response)
	err := c.cc.Invoke(ctx, "/proto.Plugin/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
type PluginServer interface {
	Info(context.Context, *empty.Empty) (*PluginInfo, error)
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	Configure(context.Context, *PluginConfig) (*empty.Empty, error)
	//Graceful Shutdown Plugin
	Stop(context.Context, *Stop_Request) (*Stop_Response, error)
}

func RegisterPluginServer(s *grpc.Server, srv PluginServer) {
	s.RegisterService(&_Plugin_serviceDesc, srv)
}

func _Plugin_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plugin/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Info(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plugin/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plugin/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Configure(ctx, req.(*PluginConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Stop_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plugin/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Stop(ctx, req.(*Stop_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Plugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _Plugin_Info_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _Plugin_Check_Handler,
		},
		{
			MethodName: "Configure",
			Handler:    _Plugin_Configure_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Plugin_Stop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/plugin.proto",
}
