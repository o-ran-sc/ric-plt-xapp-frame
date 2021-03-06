// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SubscriptionInstance subscription instance
//
// swagger:model SubscriptionInstance
type SubscriptionInstance struct {

	// e2 event instance Id
	// Required: true
	// Maximum: 65535
	// Minimum: 0
	E2EventInstanceID *int64 `json:"E2EventInstanceId"`

	// Empty string when no error.
	// Required: true
	ErrorCause *string `json:"ErrorCause"`

	// xapp event instance Id
	// Required: true
	// Maximum: 65535
	// Minimum: 0
	XappEventInstanceID *int64 `json:"XappEventInstanceId"`
}

// Validate validates this subscription instance
func (m *SubscriptionInstance) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateE2EventInstanceID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateErrorCause(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateXappEventInstanceID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SubscriptionInstance) validateE2EventInstanceID(formats strfmt.Registry) error {

	if err := validate.Required("E2EventInstanceId", "body", m.E2EventInstanceID); err != nil {
		return err
	}

	if err := validate.MinimumInt("E2EventInstanceId", "body", int64(*m.E2EventInstanceID), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("E2EventInstanceId", "body", int64(*m.E2EventInstanceID), 65535, false); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionInstance) validateErrorCause(formats strfmt.Registry) error {

	if err := validate.Required("ErrorCause", "body", m.ErrorCause); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionInstance) validateXappEventInstanceID(formats strfmt.Registry) error {

	if err := validate.Required("XappEventInstanceId", "body", m.XappEventInstanceID); err != nil {
		return err
	}

	if err := validate.MinimumInt("XappEventInstanceId", "body", int64(*m.XappEventInstanceID), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("XappEventInstanceId", "body", int64(*m.XappEventInstanceID), 65535, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionInstance) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionInstance) UnmarshalBinary(b []byte) error {
	var res SubscriptionInstance
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
