package plugin

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/pkg/errors"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type Config struct {
	ProviderType CloudProviderType `mapstructure:"cloud_provider_type"`
	// TODO: refactor config when multiple providers become being supported
	AWSAccessKeyID     string        `mapstructure:"aws_access_key_id"`
	AWSSecretAccessKey string        `mapstructure:"aws_secret_access_key"`
	AWSRegion          string        `mapstructure:"aws_region"`
	CheckInterval      time.Duration `mapstructure:"check_interval"`
	CheckTimeout       time.Duration `mapstructure:"check_timeout"`
}

type PluginsSet map[string]proto.PluginClient

type CloudProviderType string

const (
	AWSCloudProviderType CloudProviderType = "aws"
	DOCloudProviderType  CloudProviderType = "do"
)

// TODO: refactor and implement real pluggability
func (ps PluginsSet) Load(plugin proto.PluginClient, cfg *proto.PluginConfig) error {
	ctx, configureCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer configureCancel()
	_, err := plugin.Configure(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "unable to configure plugin")
	}

	ctx, infoCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer infoCancel()
	info, err := plugin.Info(ctx, &empty.Empty{})
	if err != nil {
		return errors.Wrap(err, "unable to get plugin info")
	}

	ps[info.Id] = plugin

	return nil
}

func (c Config) Validate() error {
	if c.CheckInterval.Nanoseconds() <= 0 {
		return errors.New("plugin check interval can't be less or equal to zero")
	}

	if c.CheckTimeout.Nanoseconds() <= 0 {
		return errors.New("plugin check timeout can't be less or equal to zero")
	}

	if c.ProviderType != AWSCloudProviderType {
		return errors.New("Currently only AWS is implemented")
	}

	if strings.TrimSpace(c.AWSAccessKeyID) == "" {
		return errors.New("AWS AccessKeyID was not configured")
	}

	if strings.TrimSpace(c.AWSSecretAccessKey) == "" {
		return errors.New("AWS SecretAccessKey was not configured")
	}

	if strings.TrimSpace(c.AWSRegion) == "" {
		return errors.New("AWS region was not configured")
	}

	return nil
}

func (c Config) ToProtoConfig() *proto.PluginConfig {
	return &proto.PluginConfig{
		ProviderType: newProviderType(c.ProviderType),
		CloudProviderConfig: &proto.PluginConfig_AwsConfig{
			AwsConfig: &proto.AwsConfig{
				AccessKeyId:     c.AWSAccessKeyID,
				SecretAccessKey: c.AWSSecretAccessKey,
				Region:          c.AWSRegion,
			},
		},
	}
}

func newProviderType(t CloudProviderType) proto.CloudProviderType {
	switch t {
	case DOCloudProviderType:
		return proto.CloudProviderType_DO
	default:
		return proto.CloudProviderType_AWS
	}
}
