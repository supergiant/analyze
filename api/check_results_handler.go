package api

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/supergiant/robot/pkg/storage"
	"github.com/supergiant/robot/swagger/gen/models"
	"github.com/supergiant/robot/swagger/gen/restapi/operations"
)

type checkResultsHandler struct {
	storage storage.Interface
}

func NewCheckResultsHandler(storage storage.Interface) operations.GetCheckResultsHandler {
	return &checkResultsHandler{
		storage: storage,
	}
}

func (h *checkResultsHandler) Handle(params operations.GetCheckResultsParams) middleware.Responder {

	resultsRaw, err := h.storage.GetAll(context.Background(), "/robot/check_results/")

	if err != nil {
		r := operations.NewGetCheckResultsDefault(500)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    500,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetCheckResultsOKBody{
		CheckResults: []*models.CheckResult{},
		TotalCount:   int64(len(resultsRaw)),
	}

	for _, rawResult := range resultsRaw {
		checkResult := &models.CheckResult{}
		err := checkResult.UnmarshalBinary(rawResult)
		if err != nil {
			r := operations.NewGetCheckResultsDefault(500)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    500,
				Message: &msg,
			}
			return r
		}
		result.CheckResults = append(result.CheckResults, checkResult)
	}

	return operations.NewGetCheckResultsOK().WithPayload(result)
}