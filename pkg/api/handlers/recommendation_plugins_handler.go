package handlers

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/supergiant/robot/pkg/api/operations"
	"github.com/supergiant/robot/pkg/models"
	"github.com/supergiant/robot/pkg/storage"
)

type recommendationPluginsHandler struct {
	storage storage.Interface
}

func NewRecommendationPluginsHandler(storage storage.Interface) operations.GetRecommendationPluginsHandler {
	return &recommendationPluginsHandler{
		storage: storage,
	}
}

func (h *recommendationPluginsHandler) Handle(params operations.GetRecommendationPluginsParams) middleware.Responder {

	pluginRaw, err := h.storage.GetAll(context.Background(), "/robot/plugins/")

	if err != nil {
		r := operations.NewGetRecommendationPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetRecommendationPluginsOKBody{
		InstalledRecommendationPlugins: []*models.RecommendationPlugin{},
	}

	for _, rawPlugin := range pluginRaw {
		p := &models.RecommendationPlugin{}
		err := p.UnmarshalBinary(rawPlugin)
		if err != nil {
			r := operations.NewGetRecommendationPluginsDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}
		result.InstalledRecommendationPlugins = append(result.InstalledRecommendationPlugins, p)
	}

	return operations.NewGetRecommendationPluginsOK().WithPayload(result)
}