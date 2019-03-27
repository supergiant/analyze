package analyze

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/proxy"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage"
)

type PluginController struct {
	pluginChangeEvents <-chan storage.WatchEvent
	stor               storage.Interface
	kubeClient         kube.Interface
	logger             logrus.FieldLogger
	scheduler          scheduler.Interface
	proxySet           *proxy.Set
	pluginClients      map[string]*plugin.Client
	stop               chan struct{}
}

func NewPluginController(
	events <-chan storage.WatchEvent,
	stor storage.Interface,
	kubeClient kube.Interface,
	scheduler scheduler.Interface,
	set *proxy.Set,
	logger logrus.FieldLogger,
) *PluginController {
	pc := &PluginController{
		pluginChangeEvents: events,
		stor:               stor,
		kubeClient:         kubeClient,
		scheduler:          scheduler,
		proxySet:           set,
		pluginClients:      make(map[string]*plugin.Client),
		logger:             logger,
		stop:               make(chan struct{}),
	}

	go pc.handlePluginChange()
	return pc
}

func (pc *PluginController) handlePluginChange() {
	for {
		select {
		case <-pc.stop:
			return
		case we := <-pc.pluginChangeEvents:
			if err := pc.handleEvent(we); err != nil {
				pc.logger.Errorf("unable to handle watch event, err: %+v", err)
			}
		}
	}
}

// TODO: maybe split config and plugin updates?
func (pc *PluginController) handleEvent(we storage.WatchEvent) error {
	if we.Type() == storage.Error {
		return errors.Wrap(we.Err(), "plugin watchEvent returned error")
	}
	var err error

	pluginEntry := &models.Plugin{}
	err = pluginEntry.UnmarshalBinary(we.Payload()) // TODO: handle DELETE logic, it requires storage interface update
	if err != nil {
		pc.logger.Errorf("we payload '%s'", string(we.Payload()))
		return errors.Errorf("unable to unmarshal pluginEntry entity err: %v", err)
	}

	if we.Type() == storage.Deleted {
		return pc.unregisterPlugin(pluginEntry)
	}

	if we.Type() == storage.Added {
		return pc.registerPlugin(pluginEntry)
	}

	if we.Type() == storage.Modified {
		err = pc.unregisterPlugin(pluginEntry)
		if err != nil {
			pc.logger.Errorf("unable to unregister modified plugin, error: %v", err)
		}
		return pc.registerPlugin(pluginEntry)
	}

	return nil
}

// newPluginClient returns plugin client instance which is configured and ready for work
func (pc *PluginController) newPluginClient(pluginEntry *models.Plugin) (*plugin.Client, error) {

	// pluginEntry.ServiceLabels
	ss := strings.Split(pluginEntry.ServiceEndpoint, ":")
	if len(ss) < 2 {
		return nil, errors.Errorf("unable get service name, %+v", pluginEntry.ServiceEndpoint)
	}
	service, err := pc.kubeClient.GetService(ss[0], pluginEntry.ServiceLabels)
	if err != nil {
		return nil, errors.Errorf("unable to find service for registered plugin, %+v", pluginEntry)
	}

	var servicePort string
	for _, port := range service.Spec.Ports {
		if port.Name == "grpc" {
			servicePort = strconv.Itoa(int(port.Port))
		}
	}

	if servicePort == "" {
		return nil, errors.Errorf(
			"unable to find service  port for plugin, %+v, list of ports: %+v",
			pluginEntry,
			service.Spec.Ports,
		)
	}

	pluginClient, err := plugin.NewClient(service.Spec.ClusterIP + ":" + servicePort)
	if err != nil {
		return nil, errors.Errorf("unable to instantiate pluginClient client for entity, %+v", pluginEntry)
	}

	return pluginClient, nil
}

func (pc *PluginController) unregisterPlugin(pluginEntry *models.Plugin) error {
	pluginClient, exists := pc.pluginClients[pluginEntry.ID]
	if !exists {
		return errors.Errorf("unable to find pluginClient, name: %v", pluginEntry.ID)
	}

	delete(pc.pluginClients, pluginEntry.ID)

	err := pc.scheduler.RemoveJob(pluginEntry.ID)
	if err != nil {
		return errors.Errorf("unable to find plugin job in scheduler, id: %v ", pluginEntry.ID)
	}

	err = pluginClient.Close()
	if err != nil {
		return errors.Errorf("unable to close pluginClient, id: %v, err: %+v", pluginEntry.ID, err)
	}

	pc.proxySet.RemoveProxy(pluginEntry.ID)

	return nil
}

func (pc *PluginController) registerPlugin(pluginEntry *models.Plugin) error {
	var pluginConfig *proto.PluginConfig
	pluginClient, err := pc.newPluginClient(pluginEntry)
	if err != nil {
		return errors.Wrap(err, "can't create plugin client for watchEvent")
	}

	// TODO: make it configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	pluginInfo, err := pluginClient.Info(ctx, &empty.Empty{})
	if err != nil {
		return errors.Errorf(
			"unable to load pluginInfo, name: %v, target: %s, error %v",
			pluginEntry.ID,
			pluginClient.GetConnTarget(),
			err,
		)
	}
	pc.logger.Infof("plugin info loaded %+v", *pluginInfo)

	rawConfig, err := pc.stor.Get(ctx, models.PluginConfigPrefix, pluginInfo.Id)
	pc.logger.Infof("plugin raw config %s, error: %+v", rawConfig, err)
	if err != nil && err != storage.ErrNotFound {
		return errors.Errorf("unable to load plugin config, id: %v, error %v", pluginEntry.ID, err)
	}

	if err == storage.ErrNotFound && pluginInfo.DefaultConfig == nil {
		return errors.Errorf("unable to register plugin, without default config, id: %v", pluginEntry.ID)
	}

	// when plugin is just installed we need to take default config from plugin and save it
	if err == storage.ErrNotFound {
		pluginConfig = pluginInfo.DefaultConfig
		pc.logger.Infof("plugin %s config not found, default config %+v", pluginInfo.Id, pluginInfo.DefaultConfig)

		var pluginSpecificConfig = make(map[string]interface{})
		marshalErr := json.Unmarshal(pluginInfo.DefaultConfig.PluginSpecificConfig, &pluginSpecificConfig)
		if marshalErr != nil {
			return errors.Errorf("unable to marshal plugin default config, id: %v, error %v", pluginEntry.ID, marshalErr)
		}
		pc.logger.Infof(
			"plugin %s specific config unmarshalled %+v",
			pluginInfo.Id,
			pluginSpecificConfig,
		)

		pluginConfigEntry := &models.PluginConfig{
			ExecutionInterval:    pluginConfig.ExecutionInterval.Seconds,
			PluginSpecificConfig: &pluginSpecificConfig,
		}

		b, marshalErr := pluginConfigEntry.MarshalBinary()
		if marshalErr != nil {
			return errors.Errorf("unable to marshal plugin default config, id: %v, error %v", pluginEntry.ID, marshalErr)
		}
		putError := pc.stor.Put(ctx, models.PluginConfigPrefix, pluginInfo.Id, msg(b))
		if putError != nil {
			return errors.Errorf("unable to store plugin default config, id: %v, error %v", pluginEntry.ID, putError)
		}
	}

	// true means plugin config has been read from storage and we need to apply it (case when analyze core was rebooted)
	if pluginConfig == nil {
		pluginConfigEntry := &models.PluginConfig{}
		err = pluginConfigEntry.UnmarshalBinary(rawConfig.Payload())
		if err != nil {
			return errors.Errorf("unable to unmarshal plugin config, id: %v, error %v", pluginEntry.ID, err)
		}

		bytes, err := json.Marshal(pluginConfigEntry.PluginSpecificConfig)
		if err != nil {
			return errors.Errorf("unable to marshal plugin config, id: %v, error %v", pluginEntry.ID, err)
		}

		// TODO: populate other properties
		pluginConfig = &proto.PluginConfig{
			ExecutionInterval:    ptypes.DurationProto(time.Second * time.Duration(pluginConfigEntry.ExecutionInterval)),
			PluginSpecificConfig: bytes,
		}
	}

	pc.pluginClients[pluginInfo.Id] = pluginClient

	_, err = pluginClient.Configure(ctx, pluginConfig)
	if err != nil {
		return errors.Wrap(err, "unable to configure plugin")
	}

	interval, err := ptypes.Duration(pluginConfig.ExecutionInterval)
	if err != nil {
		return errors.Wrap(err, "unable to parse execution interval for plugin")
	}

	err = pc.scheduler.ScheduleJob(pluginInfo.Id, interval, pc.check(pluginInfo.Id, pluginClient))
	if err != nil {
		return errors.Wrap(err, "unable to schedule job for plugin")
	}

	err = pc.proxySet.SetProxy(pluginEntry)
	if err != nil {
		return errors.Wrap(err, "unable to register proxy for plugin")
	}

	return nil
}

func (pc *PluginController) check(pluginID string, pluginClient *plugin.Client) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()
		checkResponse, err := pluginClient.Check(ctx, &proto.CheckRequest{})
		if err != nil {

			return errors.Errorf("unable to execute check for pluginClient: %s, error: %v", pluginID, err)
		}
		if checkResponse.Error != "" {

			return errors.Errorf("pluginClient: %s, returned error: %s", pluginID, checkResponse.Error)
		}
		if checkResponse.Result == nil {

			return errors.Errorf("pluginClient: %s, returned nil Result", pluginID)
		}

		r := checkResponse.Result

		var currentTime = time.Now()
		checkResult := models.CheckResult{
			CheckStatus:     r.GetStatus().String(),
			CompletedAt:     strfmt.DateTime(currentTime),
			Description:     string(r.GetDescription().Value),
			ExecutionStatus: r.GetExecutionStatus(),
			ID:              r.GetName(),
			Name:            r.GetName(),
		}

		bytes, err := checkResult.MarshalBinary()
		if err != nil {

			return errors.Errorf("unable to marshal check result, pluginClient: %s, returned error: %s", pluginID, err)
		}

		err = pc.stor.Put(ctx, models.CheckResultPrefix, pluginID, msg(bytes))
		if err != nil {

			return errors.Errorf("unable to store check result, pluginClient: %s, returned error: %s", pluginID, err)
		}

		return nil
	}
}

func (pc *PluginController) Stop() {
	pc.stop <- struct{}{}
}

type msg []byte

func (d msg) Payload() []byte {
	return d
}
