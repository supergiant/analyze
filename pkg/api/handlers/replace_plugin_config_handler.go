package handlers

import (
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

type replacePluginConfigHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewReplacePluginConfigHandler(storage storage.Interface, logger logrus.FieldLogger) operations.ReplacePluginConfigHandler {
	return &replacePluginConfigHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *replacePluginConfigHandler) Handle(params operations.ReplacePluginConfigParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	defer h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	r := operations.NewGetPluginConfigDefault(http.StatusInternalServerError)
	msg := "not implemented"
	r.Payload = &models.Error{
		Code:    http.StatusInternalServerError,
		Message: &msg,
	}
	return r
}
