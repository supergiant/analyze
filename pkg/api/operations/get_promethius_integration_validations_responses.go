// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"./pkg/models"
)

// GetPromethiusIntegrationValidationsOKCode is the HTTP code returned for type GetPromethiusIntegrationValidationsOK
const GetPromethiusIntegrationValidationsOKCode int = 200

/*GetPromethiusIntegrationValidationsOK no error

swagger:response getPromethiusIntegrationValidationsOK
*/
type GetPromethiusIntegrationValidationsOK struct {

	/*
	  In: Body
	*/
	Payload []*models.IntegrationComponent `json:"body,omitempty"`
}

// NewGetPromethiusIntegrationValidationsOK creates GetPromethiusIntegrationValidationsOK with default headers values
func NewGetPromethiusIntegrationValidationsOK() *GetPromethiusIntegrationValidationsOK {

	return &GetPromethiusIntegrationValidationsOK{}
}

// WithPayload adds the payload to the get promethius integration validations o k response
func (o *GetPromethiusIntegrationValidationsOK) WithPayload(payload []*models.IntegrationComponent) *GetPromethiusIntegrationValidationsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get promethius integration validations o k response
func (o *GetPromethiusIntegrationValidationsOK) SetPayload(payload []*models.IntegrationComponent) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPromethiusIntegrationValidationsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*models.IntegrationComponent, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetPromethiusIntegrationValidationsDefault error

swagger:response getPromethiusIntegrationValidationsDefault
*/
type GetPromethiusIntegrationValidationsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetPromethiusIntegrationValidationsDefault creates GetPromethiusIntegrationValidationsDefault with default headers values
func NewGetPromethiusIntegrationValidationsDefault(code int) *GetPromethiusIntegrationValidationsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetPromethiusIntegrationValidationsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get promethius integration validations default response
func (o *GetPromethiusIntegrationValidationsDefault) WithStatusCode(code int) *GetPromethiusIntegrationValidationsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get promethius integration validations default response
func (o *GetPromethiusIntegrationValidationsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get promethius integration validations default response
func (o *GetPromethiusIntegrationValidationsDefault) WithPayload(payload *models.Error) *GetPromethiusIntegrationValidationsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get promethius integration validations default response
func (o *GetPromethiusIntegrationValidationsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPromethiusIntegrationValidationsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
