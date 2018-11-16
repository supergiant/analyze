// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Plugin plugin represents the installed recommendation plugin
// swagger:model plugin
type Plugin struct {

	// detailed plugin description
	Description string `json:"description,omitempty"`

	// unique ID of installed plugin
	// basically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check
	//
	ID string `json:"id,omitempty"`

	// date/Time the plugin was installed
	// Format: date-time
	InstalledAt strfmt.DateTime `json:"installedAt,omitempty"`

	// name is the name of the plugin.
	Name string `json:"name,omitempty"`

	// plugin status
	Status string `json:"status,omitempty"`

	// plugin version, major version shall be equal to robots version
	Version string `json:"version,omitempty"`
}

// Validate validates this plugin
func (m *Plugin) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInstalledAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Plugin) validateInstalledAt(formats strfmt.Registry) error {

	if swag.IsZero(m.InstalledAt) { // not required
		return nil
	}

	if err := validate.FormatOf("installedAt", "body", "date-time", m.InstalledAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Plugin) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Plugin) UnmarshalBinary(b []byte) error {
	var res Plugin
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}