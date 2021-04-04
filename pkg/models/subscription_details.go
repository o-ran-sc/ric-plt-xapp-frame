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

// SubscriptionDetails subscription details
// swagger:model SubscriptionDetails
type SubscriptionDetails struct {

	// action to be setup list
	// Required: true
	ActionToBeSetupList ActionsToBeSetup `json:"ActionToBeSetupList"`

	// event trigger list
	// Required: true
	EventTriggerList *EventTriggerDefinition `json:"EventTriggerList"`
}

// Validate validates this subscription details
func (m *SubscriptionDetails) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActionToBeSetupList(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEventTriggerList(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SubscriptionDetails) validateActionToBeSetupList(formats strfmt.Registry) error {

	if err := validate.Required("ActionToBeSetupList", "body", m.ActionToBeSetupList); err != nil {
		return err
	}

	if err := m.ActionToBeSetupList.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("ActionToBeSetupList")
		}
		return err
	}

	return nil
}

func (m *SubscriptionDetails) validateEventTriggerList(formats strfmt.Registry) error {

	if err := validate.Required("EventTriggerList", "body", m.EventTriggerList); err != nil {
		return err
	}

	if m.EventTriggerList != nil {
		if err := m.EventTriggerList.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("EventTriggerList")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionDetails) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionDetails) UnmarshalBinary(b []byte) error {
	var res SubscriptionDetails
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
