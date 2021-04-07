// Code generated by go-swagger; DO NOT EDIT.

package clientmodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SubscriptionParams subscription params
// swagger:model SubscriptionParams
type SubscriptionParams struct {

	// client endpoint
	// Required: true
	ClientEndpoint *SubscriptionParamsClientEndpoint `json:"ClientEndpoint"`

	// meid
	// Required: true
	Meid *string `json:"Meid"`

	// r a n function ID
	// Required: true
	// Maximum: 4095
	// Minimum: 0
	RANFunctionID *int64 `json:"RANFunctionID"`

	// subscription details
	// Required: true
	SubscriptionDetails SubscriptionDetailsList `json:"SubscriptionDetails"`
}

// Validate validates this subscription params
func (m *SubscriptionParams) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClientEndpoint(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMeid(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRANFunctionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubscriptionDetails(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SubscriptionParams) validateClientEndpoint(formats strfmt.Registry) error {

	if err := validate.Required("ClientEndpoint", "body", m.ClientEndpoint); err != nil {
		return err
	}

	if m.ClientEndpoint != nil {
		if err := m.ClientEndpoint.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ClientEndpoint")
			}
			return err
		}
	}

	return nil
}

func (m *SubscriptionParams) validateMeid(formats strfmt.Registry) error {

	if err := validate.Required("Meid", "body", m.Meid); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionParams) validateRANFunctionID(formats strfmt.Registry) error {

	if err := validate.Required("RANFunctionID", "body", m.RANFunctionID); err != nil {
		return err
	}

	if err := validate.MinimumInt("RANFunctionID", "body", int64(*m.RANFunctionID), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("RANFunctionID", "body", int64(*m.RANFunctionID), 4095, false); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionParams) validateSubscriptionDetails(formats strfmt.Registry) error {

	if err := validate.Required("SubscriptionDetails", "body", m.SubscriptionDetails); err != nil {
		return err
	}

	if err := m.SubscriptionDetails.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("SubscriptionDetails")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionParams) UnmarshalBinary(b []byte) error {
	var res SubscriptionParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// SubscriptionParamsClientEndpoint xApp service address and port
// swagger:model SubscriptionParamsClientEndpoint
type SubscriptionParamsClientEndpoint struct {

	// xApp HTTP service address port
	// Maximum: 65535
	// Minimum: 0
	HTTPPort *int64 `json:"HTTPPort,omitempty"`

	// xApp service address name like 'service-ricxapp-xappname-http.ricxapp'
	Host string `json:"Host,omitempty"`

	// xApp RMR service address port
	// Maximum: 65535
	// Minimum: 0
	RMRPort *int64 `json:"RMRPort,omitempty"`
}

// Validate validates this subscription params client endpoint
func (m *SubscriptionParamsClientEndpoint) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHTTPPort(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRMRPort(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SubscriptionParamsClientEndpoint) validateHTTPPort(formats strfmt.Registry) error {

	if swag.IsZero(m.HTTPPort) { // not required
		return nil
	}

	if err := validate.MinimumInt("ClientEndpoint"+"."+"HTTPPort", "body", int64(*m.HTTPPort), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("ClientEndpoint"+"."+"HTTPPort", "body", int64(*m.HTTPPort), 65535, false); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionParamsClientEndpoint) validateRMRPort(formats strfmt.Registry) error {

	if swag.IsZero(m.RMRPort) { // not required
		return nil
	}

	if err := validate.MinimumInt("ClientEndpoint"+"."+"RMRPort", "body", int64(*m.RMRPort), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("ClientEndpoint"+"."+"RMRPort", "body", int64(*m.RMRPort), 65535, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionParamsClientEndpoint) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionParamsClientEndpoint) UnmarshalBinary(b []byte) error {
	var res SubscriptionParamsClientEndpoint
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
