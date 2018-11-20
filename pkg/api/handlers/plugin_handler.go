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

type pluginsHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewPluginsHandler(storage storage.Interface, logger logrus.FieldLogger) operations.GetPluginsHandler {
	return &pluginsHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *pluginsHandler) Handle(params operations.GetPluginsParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	pluginRaw, err := h.storage.GetAll(ctx, models.PluginPrefix)

	if err != nil {
		r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetPluginsOKBody{
		InstalledPlugins: []*models.Plugin{},
	}

	for _, rawPlugin := range pluginRaw {
		p := &models.Plugin{}
		err := p.UnmarshalBinary(rawPlugin)
		if err != nil {
			r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}
		result.InstalledPlugins = append(result.InstalledPlugins, p)
	}
	h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	return operations.NewGetPluginsOK().WithPayload(result)
}
