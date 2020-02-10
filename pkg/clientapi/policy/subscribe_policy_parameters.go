// Code generated by go-swagger; DO NOT EDIT.

package policy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	clientmodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
)

// NewSubscribePolicyParams creates a new SubscribePolicyParams object
// with the default values initialized.
func NewSubscribePolicyParams() *SubscribePolicyParams {
	var ()
	return &SubscribePolicyParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSubscribePolicyParamsWithTimeout creates a new SubscribePolicyParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSubscribePolicyParamsWithTimeout(timeout time.Duration) *SubscribePolicyParams {
	var ()
	return &SubscribePolicyParams{

		timeout: timeout,
	}
}

// NewSubscribePolicyParamsWithContext creates a new SubscribePolicyParams object
// with the default values initialized, and the ability to set a context for a request
func NewSubscribePolicyParamsWithContext(ctx context.Context) *SubscribePolicyParams {
	var ()
	return &SubscribePolicyParams{

		Context: ctx,
	}
}

// NewSubscribePolicyParamsWithHTTPClient creates a new SubscribePolicyParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSubscribePolicyParamsWithHTTPClient(client *http.Client) *SubscribePolicyParams {
	var ()
	return &SubscribePolicyParams{
		HTTPClient: client,
	}
}

/*SubscribePolicyParams contains all the parameters to send to the API endpoint
for the subscribe policy operation typically these are written to a http.Request
*/
type SubscribePolicyParams struct {

	/*PolicyParams
	  Subscription policy parameters

	*/
	PolicyParams *clientmodel.PolicyParams

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the subscribe policy params
func (o *SubscribePolicyParams) WithTimeout(timeout time.Duration) *SubscribePolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the subscribe policy params
func (o *SubscribePolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the subscribe policy params
func (o *SubscribePolicyParams) WithContext(ctx context.Context) *SubscribePolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the subscribe policy params
func (o *SubscribePolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the subscribe policy params
func (o *SubscribePolicyParams) WithHTTPClient(client *http.Client) *SubscribePolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the subscribe policy params
func (o *SubscribePolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPolicyParams adds the policyParams to the subscribe policy params
func (o *SubscribePolicyParams) WithPolicyParams(policyParams *clientmodel.PolicyParams) *SubscribePolicyParams {
	o.SetPolicyParams(policyParams)
	return o
}

// SetPolicyParams adds the policyParams to the subscribe policy params
func (o *SubscribePolicyParams) SetPolicyParams(policyParams *clientmodel.PolicyParams) {
	o.PolicyParams = policyParams
}

// WriteToRequest writes these params to a swagger request
func (o *SubscribePolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PolicyParams != nil {
		if err := r.SetBodyParam(o.PolicyParams); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
