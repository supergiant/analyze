package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/supergiant/analyze/pkg/config"

	"github.com/spf13/viper"
)

type mainConfig struct {
	Connection connectionConfig `mapstructure:"connection"`
	Pods       []string         `mapstructure:"pods"`
}

type connectionConfig struct {
	Vhost                       string        `mapstructure:"vhost"`
	Server                      string        `mapstructure:"server"`
	Port                        string        `mapstructure:"port"`
	Username                    string        `mapstructure:"username"`
	Password                    string        `mapstructure:"password"`
	UseTLS                      bool          `mapstructure:"use_tls"`
	SkipTLSCertificateVerifying bool          `mapstructure:"skip_tls_certificate_verifying"`
	ConnectionTries             int           `mapstructure:"connection_tries"`
	ConnectionTryTimeout        time.Duration `mapstructure:"connection_try_timeout"`
}

func getRefConf() mainConfig {
	return mainConfig{
		Connection: connectionConfig{
			Vhost:                       "/test/123",
			Server:                      "rpc.qbox.io",
			Port:                        "5671",
			Username:                    "test",
			Password:                    "",
			UseTLS:                      true,
			SkipTLSCertificateVerifying: true,
			ConnectionTries:             3,
			ConnectionTryTimeout:        5 * time.Second,
		},
		Pods: []string{"pod1", "pod2", "pod3"},
	}
}

func TestConfig_RefreshFromFileSuccessfully(t *testing.T) {
	viper.Reset()
	c := &mainConfig{}
	err := config.ReadFromFiles(c, []string{"./test-fixtures/fixture-config.yml"})
	if err != nil {
		t.Fatalf("unable to get config from file %+v", err)
	}

	if getRefConf().Connection != c.Connection {
		t.Fatal("config loaded incorrectly")
	}

	for i := range getRefConf().Pods {
		if getRefConf().Pods[i] != c.Pods[i] {
			t.Fatal("config loaded incorrectly")
		}
	}
}

func TestConfig_MergeEnvSuccessfully(t *testing.T) {
	viper.Reset()
	c := &mainConfig{}
	err := config.ReadFromFiles(c, []string{"./test-fixtures/fixture-config.yml"})
	if err != nil {
		t.Fatalf("unable to get config from file %+v", err)
	}

	os.Setenv("RK_PODS", "pod1,pod2,pod3,pod4")
	err = config.MergeEnv("RK", c)
	if err != nil {
		t.Fatalf("unable override config properties from env %+v", err)
	}

	tr := []string{"pod1", "pod2", "pod3", "pod4"}
	for i := range c.Pods {
		if c.Pods[i] != tr[i] {
			t.Fatal("config loaded incorrectly")
		}
	}
}

func TestConfig_MergeEnvSuccessfullyWhenEnvSetButEmpty(t *testing.T) {
	viper.Reset()
	c := &mainConfig{}
	err := config.ReadFromFiles(c, []string{"./test-fixtures/fixture-config.yml"})
	if err != nil {
		t.Fatalf("unable to get config from file %+v", err)
	}

	os.Setenv("RK_PODS", "")
	err = config.MergeEnv("RK", c)
	if err != nil {
		t.Fatalf("unable override config properties from env %+v", err)
	}

	tr := []string{"pod1", "pod2", "pod3"}
	for i := range c.Pods {
		if c.Pods[i] != tr[i] {
			t.Fatal("config loaded incorrectly")
		}
	}
}

func TestConfig_MergeEnvSuccessfullyWhenEnvSetPartially(t *testing.T) {
	viper.Reset()
	c := &mainConfig{}
	err := config.ReadFromFiles(c, []string{"./test-fixtures/fixture-config.yml"})
	if err != nil {
		t.Fatalf("unable to get config from file %+v", err)
	}

	os.Setenv("RK_PODS", "pod33")
	err = config.MergeEnv("RK", c)
	if err != nil {
		t.Fatalf("unable override config properties from env %+v", err)
	}

	tr := []string{"pod33", "pod2", "pod3"}
	for i := range c.Pods {
		if c.Pods[i] != tr[i] {
			t.Fatal("config loaded incorrectly")
		}
	}
}
