// Code generated by go-swagger; DO NOT EDIT.

package control

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	models "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
)

// NewSubscribeControlParams creates a new SubscribeControlParams object
// no default values defined in spec.
func NewSubscribeControlParams() SubscribeControlParams {

	return SubscribeControlParams{}
}

// SubscribeControlParams contains all the bound params for the subscribe control operation
// typically these are obtained from a http.Request
//
// swagger:parameters subscribeControl
type SubscribeControlParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Subscription control parameters
	  In: body
	*/
	ControlParams *models.ControlParams
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSubscribeControlParams() beforehand.
func (o *SubscribeControlParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.ControlParams
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("controlParams", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.ControlParams = &body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
