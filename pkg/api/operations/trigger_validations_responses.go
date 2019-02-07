// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/supergiant/analyze/pkg/models"
)

// TriggerValidationsNoContentCode is the HTTP code returned for type TriggerValidationsNoContent
const TriggerValidationsNoContentCode int = 204

/*TriggerValidationsNoContent validation has been triggered

swagger:response triggerValidationsNoContent
*/
type TriggerValidationsNoContent struct {
}

// NewTriggerValidationsNoContent creates TriggerValidationsNoContent with default headers values
func NewTriggerValidationsNoContent() *TriggerValidationsNoContent {

	return &TriggerValidationsNoContent{}
}

// WriteResponse to the client
func (o *TriggerValidationsNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*TriggerValidationsDefault error

swagger:response triggerValidationsDefault
*/
type TriggerValidationsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewTriggerValidationsDefault creates TriggerValidationsDefault with default headers values
func NewTriggerValidationsDefault(code int) *TriggerValidationsDefault {
	if code <= 0 {
		code = 500
	}

	return &TriggerValidationsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the trigger validations default response
func (o *TriggerValidationsDefault) WithStatusCode(code int) *TriggerValidationsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the trigger validations default response
func (o *TriggerValidationsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the trigger validations default response
func (o *TriggerValidationsDefault) WithPayload(payload *models.Error) *TriggerValidationsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the trigger validations default response
func (o *TriggerValidationsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TriggerValidationsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
