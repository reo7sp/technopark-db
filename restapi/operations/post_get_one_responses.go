// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/reo7sp/technopark-db/models"
)

// PostGetOneOKCode is the HTTP code returned for type PostGetOneOK
const PostGetOneOKCode int = 200

/*PostGetOneOK Информация о ветке обсуждения.


swagger:response postGetOneOK
*/
type PostGetOneOK struct {

	/*
	  In: Body
	*/
	Payload *models.PostFull `json:"body,omitempty"`
}

// NewPostGetOneOK creates PostGetOneOK with default headers values
func NewPostGetOneOK() *PostGetOneOK {
	return &PostGetOneOK{}
}

// WithPayload adds the payload to the post get one o k response
func (o *PostGetOneOK) WithPayload(payload *models.PostFull) *PostGetOneOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post get one o k response
func (o *PostGetOneOK) SetPayload(payload *models.PostFull) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostGetOneOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostGetOneNotFoundCode is the HTTP code returned for type PostGetOneNotFound
const PostGetOneNotFoundCode int = 404

/*PostGetOneNotFound Ветка обсуждения отсутсвует в форуме.


swagger:response postGetOneNotFound
*/
type PostGetOneNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostGetOneNotFound creates PostGetOneNotFound with default headers values
func NewPostGetOneNotFound() *PostGetOneNotFound {
	return &PostGetOneNotFound{}
}

// WithPayload adds the payload to the post get one not found response
func (o *PostGetOneNotFound) WithPayload(payload *models.Error) *PostGetOneNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post get one not found response
func (o *PostGetOneNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostGetOneNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
