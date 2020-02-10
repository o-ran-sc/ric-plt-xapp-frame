// Code generated by go-swagger; DO NOT EDIT.

package common

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// UnsubscribeReader is a Reader for the Unsubscribe structure.
type UnsubscribeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UnsubscribeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewUnsubscribeNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewUnsubscribeBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewUnsubscribeInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUnsubscribeNoContent creates a UnsubscribeNoContent with default headers values
func NewUnsubscribeNoContent() *UnsubscribeNoContent {
	return &UnsubscribeNoContent{}
}

/*UnsubscribeNoContent handles this case with default header values.

Operation done successfully
*/
type UnsubscribeNoContent struct {
}

func (o *UnsubscribeNoContent) Error() string {
	return fmt.Sprintf("[DELETE /subscriptions/{subscriptionId}][%d] unsubscribeNoContent ", 204)
}

func (o *UnsubscribeNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUnsubscribeBadRequest creates a UnsubscribeBadRequest with default headers values
func NewUnsubscribeBadRequest() *UnsubscribeBadRequest {
	return &UnsubscribeBadRequest{}
}

/*UnsubscribeBadRequest handles this case with default header values.

Invalid requestorId supplied
*/
type UnsubscribeBadRequest struct {
}

func (o *UnsubscribeBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /subscriptions/{subscriptionId}][%d] unsubscribeBadRequest ", 400)
}

func (o *UnsubscribeBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUnsubscribeInternalServerError creates a UnsubscribeInternalServerError with default headers values
func NewUnsubscribeInternalServerError() *UnsubscribeInternalServerError {
	return &UnsubscribeInternalServerError{}
}

/*UnsubscribeInternalServerError handles this case with default header values.

Internal error
*/
type UnsubscribeInternalServerError struct {
}

func (o *UnsubscribeInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /subscriptions/{subscriptionId}][%d] unsubscribeInternalServerError ", 500)
}

func (o *UnsubscribeInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
