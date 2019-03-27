package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

type pluginConfigHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewPluginConfigHandler(storage storage.Interface, logger logrus.FieldLogger) operations.GetPluginConfigHandler {
	return &pluginConfigHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *pluginConfigHandler) Handle(params operations.GetPluginConfigParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	defer h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rawConfig, err := h.storage.Get(ctx, models.PluginConfigPrefix, params.PluginID)
	if err == storage.ErrNotFound {
		return operations.NewGetPluginConfigNotFound()
	}

	if err != nil {
		r := operations.NewGetPluginConfigDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	pluginConfigEntry := &models.PluginConfig{}
	err = pluginConfigEntry.UnmarshalBinary(rawConfig.Payload())
	if err != nil {
		r := operations.NewGetPluginConfigDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	return operations.NewGetPluginConfigOK().WithPayload(pluginConfigEntry)
}
