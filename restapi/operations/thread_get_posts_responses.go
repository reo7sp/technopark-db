// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/reo7sp/technopark-db/models"
)

// ThreadGetPostsOKCode is the HTTP code returned for type ThreadGetPostsOK
const ThreadGetPostsOKCode int = 200

/*ThreadGetPostsOK Информация о сообщениях форума.


swagger:response threadGetPostsOK
*/
type ThreadGetPostsOK struct {

	/*
	  In: Body
	*/
	Payload models.Posts `json:"body,omitempty"`
}

// NewThreadGetPostsOK creates ThreadGetPostsOK with default headers values
func NewThreadGetPostsOK() *ThreadGetPostsOK {
	return &ThreadGetPostsOK{}
}

// WithPayload adds the payload to the thread get posts o k response
func (o *ThreadGetPostsOK) WithPayload(payload models.Posts) *ThreadGetPostsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thread get posts o k response
func (o *ThreadGetPostsOK) SetPayload(payload models.Posts) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThreadGetPostsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make(models.Posts, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// ThreadGetPostsNotFoundCode is the HTTP code returned for type ThreadGetPostsNotFound
const ThreadGetPostsNotFoundCode int = 404

/*ThreadGetPostsNotFound Ветка обсуждения отсутсвует в форуме.


swagger:response threadGetPostsNotFound
*/
type ThreadGetPostsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewThreadGetPostsNotFound creates ThreadGetPostsNotFound with default headers values
func NewThreadGetPostsNotFound() *ThreadGetPostsNotFound {
	return &ThreadGetPostsNotFound{}
}

// WithPayload adds the payload to the thread get posts not found response
func (o *ThreadGetPostsNotFound) WithPayload(payload *models.Error) *ThreadGetPostsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thread get posts not found response
func (o *ThreadGetPostsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThreadGetPostsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
