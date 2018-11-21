// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/supergiant/analyze/pkg/models"
)

// GetPluginsOKCode is the HTTP code returned for type GetPluginsOK
const GetPluginsOKCode int = 200

/*GetPluginsOK no error

swagger:response getPluginsOK
*/
type GetPluginsOK struct {

	/*installed plugins
	  In: Body
	*/
	Payload []*models.Plugin `json:"body,omitempty"`
}

// NewGetPluginsOK creates GetPluginsOK with default headers values
func NewGetPluginsOK() *GetPluginsOK {

	return &GetPluginsOK{}
}

// WithPayload adds the payload to the get plugins o k response
func (o *GetPluginsOK) WithPayload(payload []*models.Plugin) *GetPluginsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get plugins o k response
func (o *GetPluginsOK) SetPayload(payload []*models.Plugin) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPluginsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*models.Plugin, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetPluginsDefault error

swagger:response getPluginsDefault
*/
type GetPluginsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetPluginsDefault creates GetPluginsDefault with default headers values
func NewGetPluginsDefault(code int) *GetPluginsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetPluginsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get plugins default response
func (o *GetPluginsDefault) WithStatusCode(code int) *GetPluginsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get plugins default response
func (o *GetPluginsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get plugins default response
func (o *GetPluginsDefault) WithPayload(payload *models.Error) *GetPluginsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get plugins default response
func (o *GetPluginsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPluginsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
